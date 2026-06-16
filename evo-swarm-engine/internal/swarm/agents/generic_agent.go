package agents

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/ai/llm"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/ai/tools"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/tts"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/swarm/registry"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

// GenericAgent é o Especialista Universal que materializa a personalidade definida no Banco de Dados
type GenericAgent struct {
	EventBus *evbus.EventBus
	DB       *gorm.DB
	Model    models.Agent
	Subject  string
}

func NewGenericAgent(bus *evbus.EventBus, db *gorm.DB, agent models.Agent) *GenericAgent {
	return &GenericAgent{
		EventBus: bus,
		DB:       db,
		Model:    agent,
		Subject:  "agent." + agent.ID.String() + ".task",
	}
}

func (a *GenericAgent) Start() error {
	log.Printf("🤖 [GenericAgent] Agente '%s' (%s) iniciado. Escutando %s...", a.Model.Name, a.Model.ID, a.Subject)

	// Registra o agente no Registry de forma dinâmica
	reg, err := registry.NewRegistry(a.EventBus)
	if err != nil {
		log.Printf("❌ [GenericAgent] Erro ao inicializar Registry (NATS JetStream offline?): %v", err)
		return err
	}
	instruction := ""
	if a.Model.Instruction != nil {
		instruction = *a.Model.Instruction
	}
	
	// TODO: No futuro as skills poderão vir do banco de dados (relacionamento agent_tools)
	reg.Register(registry.Capability{
		AgentID:     a.Model.ID.String(),
		Subject:     a.Subject,
		Description: a.Model.Name + " - " + instruction,
		Skills:      []string{"dynamic_agent"},
	})

	_, err = a.EventBus.Conn.QueueSubscribe(a.Subject, "generic_agent_group", a.handleTask)
	return err
}

func (a *GenericAgent) handleTask(msg *nats.Msg) {
	var event events.AgentStartedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("❌ [GenericAgent] Erro ao decodificar AgentStartedEvent: %v", err)
		return
	}

	log.Printf("🛠️ [GenericAgent] [%s] Processando Tarefa: [%s]", a.Model.Name, event.TaskID)

	var incomingData struct {
		Source         string `json:"source"`
		Sender         string `json:"sender"`
		ContextID      string `json:"context_id"`
		Content        string `json:"content"`
		ConversationID string `json:"conversation_id"`
		AccountID      int64  `json:"account_id"`
	}
	if err := json.Unmarshal([]byte(event.Payload), &incomingData); err != nil {
		log.Printf("❌ [GenericAgent] [%s] Erro ao decodificar payload interno: %v", a.Model.Name, err)
		return
	}

	log.Printf("📋 [GenericAgent] [%s] Payload recebido: source='%s', sender='%s', content='%s', conversation_id='%s'",
		a.Model.Name, incomingData.Source, incomingData.Sender, incomingData.Content, incomingData.ConversationID)

	conversationID := incomingData.ContextID
	if conversationID == "" {
		conversationID = incomingData.Sender
	}

	// Consulta Memória de Curto Prazo e RAG (Phase 3)
	historyRequest, _ := json.Marshal(struct {
		ConversationID string `json:"conversation_id"`
		Limit          int    `json:"limit"`
	}{
		ConversationID: conversationID,
		Limit:          10,
	})

	var history []models.LLMMessage
	memReply, err := a.EventBus.Conn.Request("memory.history.get", historyRequest, 2*time.Second)
	if err == nil {
		var memResponse struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if json.Unmarshal(memReply.Data, &memResponse) == nil {
			for _, m := range memResponse.Messages {
				history = append(history, models.LLMMessage{Role: m.Role, Content: m.Content})
			}
		}
	}

	// Determina API Key e Model Dinamicamente baseados na configuração DO AGENTE
	apiKey := "sk-mock-key"
	modelName := "gpt-4o-mini"
	systemPrompt := "Você é o agente genérico. Por favor, ajude o usuário."

	if a.Model.Model != nil && *a.Model.Model != "" {
		modelName = *a.Model.Model
	}
	if a.Model.Instruction != nil && *a.Model.Instruction != "" {
		systemPrompt = *a.Model.Instruction
	}
	if a.Model.APIKeyID != nil {
		var dbKey models.APIKey
		if err := a.DB.First(&dbKey, "id = ?", a.Model.APIKeyID).Error; err == nil {
			apiKey = dbKey.Key
		}
	} else {
		var dbKey models.APIKey
		if err := a.DB.Where("is_active = ? AND provider IN ?", true, []string{"openai", "anthropic", "openrouter", "OpenRouter"}).First(&dbKey).Error; err == nil {
			apiKey = dbKey.Key
		}
	}

	// RAG Service Lookup
	ragRequest, _ := json.Marshal(struct {
		Query   string `json:"query"`
		Limit   int    `json:"limit"`
		AgentID string `json:"agent_id"`
		APIKey  string `json:"api_key"`
	}{
		Query:   incomingData.Content,
		Limit:   3,
		AgentID: a.Model.ID.String(),
		APIKey:  apiKey,
	})

	ragContext := ""
	ragReply, err := a.EventBus.Conn.Request("service.rag.query", ragRequest, 5*time.Second)
	if err == nil {
		var ragResponse struct {
			Status  string   `json:"status"`
			Context []string `json:"context"`
		}
		if err := json.Unmarshal(ragReply.Data, &ragResponse); err == nil {
			if ragResponse.Status == "success" {
				for _, ctx := range ragResponse.Context {
					ragContext += ctx + "\n---\n"
				}
				log.Printf("📚 [GenericAgent] [%s] RAG Context Injetado com %d chunks.", a.Model.Name, len(ragResponse.Context))
			}
		}
	}

	llmProvider, err := llm.NewLLMProvider(modelName, apiKey)
	if err != nil {
		log.Printf("❌ [GenericAgent] [%s] Erro ao inicializar LLM: %v", a.Model.Name, err)
		return
	}

	systemInstruction := systemPrompt
	if ragContext != "" {
		systemInstruction += "\n\n<knowledge_context>\n" + ragContext + "\n</knowledge_context>"
	}

	// Consulta o Registry e injeta Handoff Tools
	reg, err := registry.NewRegistry(a.EventBus)
	var swarmTools []models.LLMTool
	var caps []registry.Capability
	if err == nil {
		caps, _ = reg.GetAllCapabilities()
		
		// Filtra as capabilities apenas para os sub-agentes configurados no banco de dados
		var filteredCaps []registry.Capability
		var allowedSubAgents []string
		
		// a.Model.Config é do tipo datatypes.JSON, precisamos decodificar
		var configMap map[string]interface{}
		if len(a.Model.Config) > 0 {
			json.Unmarshal(a.Model.Config, &configMap)
		}
		
		if subAgentsRaw, ok := configMap["sub_agents"]; ok {
			if subAgentsList, ok := subAgentsRaw.([]interface{}); ok {
				for _, sa := range subAgentsList {
					if idStr, ok := sa.(string); ok {
						allowedSubAgents = append(allowedSubAgents, idStr)
					}
				}
			}
		}

		if len(allowedSubAgents) > 0 {
			for _, cap := range caps {
				for _, allowed := range allowedSubAgents {
					if cap.AgentID == allowed {
						filteredCaps = append(filteredCaps, cap)
						break
					}
				}
			}
		} else {
			// Se não houver sub_agents configurados, não permite handoff (comportamento estrito hierárquico)
			filteredCaps = []registry.Capability{}
		}

		swarmTools = tools.GetSwarmTools(filteredCaps, a.Model.ID.String())
		
		if len(swarmTools) > 0 {
			systemInstruction += "\n\nVocê tem acesso a outros sub-agentes especializados. Se o usuário pedir algo fora de sua área, use a ferramenta 'delegate_to_agent' para transferir a conversa."
		}
	}

	// TODO: Aqui integraríamos com o banco de ferramentas (tools) conectadas ao Agente.
	// Por enquanto, usamos apenas as de Swarm.
	agentTools := []models.LLMTool{}
	agentTools = append(agentTools, swarmTools...)

	req := models.LLMRequest{
		Model:       modelName,
		System:      systemInstruction,
		Temperature: 0.5,
		Messages:    history, 
		Tools:       agentTools,
	}
	
	req.Messages = append(req.Messages, models.LLMMessage{Role: "user", Content: incomingData.Content})

	for loop := 0; loop < 3; loop++ {
		resp, err := llmProvider.Generate(context.Background(), req)
		if err != nil {
			log.Printf("❌ [GenericAgent] [%s] Erro na LLM: %v", a.Model.Name, err)
			return
		}

		if len(resp.ToolCalls) > 0 {
			req.Messages = append(req.Messages, models.LLMMessage{
				Role:      "assistant",
				ToolCalls: resp.ToolCalls,
			})

			for _, tc := range resp.ToolCalls {
				log.Printf("🔧 [GenericAgent] [%s] Emitindo Comando via NATS: %s", a.Model.Name, tc.FunctionName)
				
				if tc.FunctionName == "delegate_to_agent" {
					var args struct {
						TargetAgent string `json:"target_agent"`
						Reason      string `json:"reason"`
					}
					json.Unmarshal([]byte(tc.Arguments), &args)

					targetSubject := ""
					for _, cap := range caps {
						if cap.AgentID == args.TargetAgent {
							targetSubject = cap.Subject
							break
						}
					}

					if targetSubject != "" {
						log.Printf("🔄 [GenericAgent] [%s] Delegando conversa para o sub-agente: %s", a.Model.Name, args.TargetAgent)
						
						handoffEvent := event
						handoffIncoming := incomingData
						handoffIncoming.Content = args.Reason
						handoffPayload, _ := json.Marshal(handoffIncoming)
						handoffEvent.Payload = string(handoffPayload)

						eventData, _ := json.Marshal(handoffEvent)
						a.EventBus.Publish(targetSubject, eventData)

						// Dispara o Webhook de Handoff para o Barramento Geral
						webhookEvent := events.AgentHandoffEvent{
							BaseEvent: events.BaseEvent{
								EventID:   event.EventID,
								EventType: events.EventAgentHandoff,
								Timestamp: time.Now(),
							},
							SourceAgentID:  a.Model.ID.String(),
							TargetAgentID:  args.TargetAgent,
							ConversationID: conversationID,
							Reason:         args.Reason,
						}
						webhookData, _ := json.Marshal(webhookEvent)
						a.EventBus.Publish(string(events.EventAgentHandoff), webhookData)

						return // Finaliza silenciosamente e deixa o sub-agente assumir
					} else {
						req.Messages = append(req.Messages, models.LLMMessage{
							Role:       "tool",
							Content:    "Erro: O agente especificado não está online ou não existe.",
							ToolCallID: tc.ID,
							Name:       tc.FunctionName,
						})
					}
				}
				// TODO: Tratar outras ferramentas dinâmicas cadastradas no banco
			}
		} else {
			log.Printf("💬 [GenericAgent] [%s] Resposta Final gerada com sucesso: '%s'", a.Model.Name, resp.Content)

			var audioURL string
			if len(a.Model.Config) > 0 {
				var configMap map[string]interface{}
				if err := json.Unmarshal(a.Model.Config, &configMap); err == nil {
					if integrations, ok := configMap["integrations"].(map[string]interface{}); ok {
						var ttsConfig map[string]interface{}
						if ttsRaw, hasTts := integrations["tts"].(map[string]interface{}); hasTts {
							ttsConfig = ttsRaw
						} else if elevenlabsRaw, hasEleven := integrations["elevenlabs"].(map[string]interface{}); hasEleven {
							ttsConfig = elevenlabsRaw
							ttsConfig["provider"] = "elevenlabs"
						}

						if ttsConfig != nil {
							providerName, _ := ttsConfig["provider"].(string)
							voiceID, _ := ttsConfig["voice_id"].(string)
							apiKey, _ := ttsConfig["api_key"].(string)

							if providerName != "" && voiceID != "" && apiKey != "" {
								log.Printf("🗣️ [GenericAgent] [%s] Gerando áudio via TTS (%s)", a.Model.Name, providerName)
								ttsFactory := tts.NewFactory()
								if provider, err := ttsFactory.GetProvider(providerName); err == nil {
									req := tts.Request{
										Text:    resp.Content,
										VoiceID: voiceID,
										APIKey:  apiKey,
									}
									if modelID, ok := ttsConfig["model_id"].(string); ok {
										req.Model = modelID
									}
									audioBytes, err := provider.GenerateSpeech(context.Background(), req)
									if err != nil {
										log.Printf("❌ [GenericAgent] Erro na geração TTS: %v", err)
									} else if len(audioBytes) > 0 {
										audioBase64 := base64.StdEncoding.EncodeToString(audioBytes)
										audioURL = fmt.Sprintf("data:audio/ogg;base64,%s", audioBase64)
										log.Printf("✅ [GenericAgent] Áudio gerado com sucesso (%d bytes)", len(audioBytes))
									}
								} else {
									log.Printf("❌ [GenericAgent] Provedor TTS não suportado: %s", providerName)
								}
							}
						}
					}
				}
			}

			type OutboundResponse struct {
				Source         string `json:"source"`
				Sender         string `json:"sender"`
				Status         string `json:"status"`
				Content        string `json:"content"`
				AudioURL       string `json:"audio_url,omitempty"`
				ConversationID string `json:"conversation_id"`
				AccountID      int64  `json:"account_id"`
			}
			responsePayload, err := json.Marshal(OutboundResponse{
				Source:         incomingData.Source,
				Sender:         incomingData.Sender,
				Status:         "completed",
				Content:        resp.Content,
				AudioURL:       audioURL,
				ConversationID: incomingData.ConversationID,
				AccountID:      incomingData.AccountID,
			})
			if err != nil {
				log.Printf("❌ [GenericAgent] [%s] Erro ao serializar resposta final: %v", a.Model.Name, err)
				return
			}
			
			log.Printf("📤 [GenericAgent] [%s] Publicando resposta em outbound.message (ConversationID: %s)", a.Model.Name, incomingData.ConversationID)
			a.EventBus.Publish("outbound.message", responsePayload)
			a.EventBus.Publish("stream."+event.TaskID.String(), responsePayload)
			break
		}
	}
}
