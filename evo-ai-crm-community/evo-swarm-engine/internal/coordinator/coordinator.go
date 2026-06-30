package coordinator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/ai/llm"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/swarm/registry"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

// Coordinator é o roteador principal do Swarm.
type Coordinator struct {
	EventBus *evbus.EventBus
	DB       *gorm.DB
}

// NewCoordinator instancia um novo Roteador.
func NewCoordinator(bus *evbus.EventBus, db *gorm.DB) *Coordinator {
	return &Coordinator{
		EventBus: bus,
		DB:       db,
	}
}

// Start inicializa o Coordinator e assina os eventos essenciais do NATS.
func (c *Coordinator) Start() error {
	log.Println("🚀 Iniciando Swarm Coordinator...")

	// Assina eventos de novos Leads (Exemplo de Ingress Point para workflows)
	_, err := c.EventBus.Conn.QueueSubscribe(string(events.EventLeadCreated), "coordinator_pool", c.handleLeadCreated)
	if err != nil {
		return err
	}

	// Assina eventos de conclusão de agentes (Handoff e Continuidade do Workflow)
	_, err = c.EventBus.Conn.QueueSubscribe(string(events.EventAgentFinished), "coordinator_pool", c.handleAgentFinished)
	if err != nil {
		return err
	}

	// Assina eventos de input do humano (Webhook WhatsApp/Chatwoot)
	_, err = c.EventBus.Conn.QueueSubscribe(string(events.EventMessageReceived), "coordinator_pool", c.handleMessageReceived)
	if err != nil {
		return err
	}

	log.Println("✅ Swarm Coordinator aguardando eventos...")
	return nil
}

// handleLeadCreated é acionado quando um Lead é criado.
// Aqui a Engine de Workflow (DAG) seria instanciada e o primeiro nó (Agente) disparado.
func (c *Coordinator) handleLeadCreated(msg *nats.Msg) {
	log.Printf("[Coordinator] Novo evento recebido: %s", msg.Subject)

	// Na prática:
	// 1. Decodifica o Lead.
	// 2. Busca no banco qual Workflow está atrelado a esse Lead/Sinal.
	// 3. Monta o DAG em memória (workflow.NewWorkflow).
	// 4. Executa o DAG passando o contexto do Lead.
}

// handleAgentFinished intercepta o fim de um trabalho de um agente.
func (c *Coordinator) handleAgentFinished(msg *nats.Msg) {
	log.Printf("[Coordinator] Agente finalizou tarefa: %s", msg.Subject)

	var event events.AgentFinishedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[Coordinator] Erro ao decodificar AgentFinishedEvent: %v", err)
		return
	}

	// Na prática:
	// 1. Busca a Engine de Workflow (DAG) pausada pelo TaskID/TraceID no Redis.
	// 2. Atualiza o estado da DAG com o Result (event.Result).
	// 3. Verifica o DAG para encontrar os próximos nós (NextNodes).
	// 4. Dispara EventAgentStarted para os próximos agentes.
}

// handleMessageReceived lida com a entrada de dados do usuário (ex: WhatsApp/Chatwoot)
func (c *Coordinator) handleMessageReceived(msg *nats.Msg) {
	log.Printf("[Coordinator] 📩 Nova mensagem recebida do usuário: %s", msg.Subject)

	var payload struct {
		Source  string `json:"source"`
		Content string `json:"content"`
		Sender  string `json:"sender"`
		AgentID string `json:"agent_id"` // Adicionado suporte para agent_id opcional
	}

	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		log.Printf("❌ [Coordinator] Erro ao decodificar payload da mensagem: %v", err)
		return
	}

	// Resolve o AgentID para vincular no Evento
	var resolvedAgentID uuid.UUID
	if payload.AgentID != "" {
		if parsed, err := uuid.Parse(payload.AgentID); err == nil {
			resolvedAgentID = parsed
		}
	}

	// Se não veio AgentID no webhook, faz o fallback para o primeiro agente LLM válido do sistema
	if resolvedAgentID == uuid.Nil {
		var defaultAgent models.Agent
		if err := c.DB.Where("type = ?", "llm").First(&defaultAgent).Error; err == nil {
			resolvedAgentID = defaultAgent.ID
			log.Printf("🔄 [Coordinator] Webhook não especificou AgentID. Fazendo fallback para o agente padrão: %s", resolvedAgentID)
		} else {
			log.Printf("⚠️ [Coordinator] Nenhum Agente padrão encontrado no banco de dados!")
		}
	}

	var targetAgent string
	var decision string

	// Carrega as capabilities ativas do Registro
	reg, err := registry.NewRegistry(c.EventBus)
	var caps []registry.Capability
	if err == nil {
		caps, _ = reg.GetAllCapabilities()
	}

	if resolvedAgentID != uuid.Nil {
		decision = resolvedAgentID.String()
		for _, cap := range caps {
			if cap.AgentID == decision {
				targetAgent = cap.Subject
				break
			}
		}
		if targetAgent == "" {
			log.Printf("⚠️ [Coordinator] Agente %s informado no payload não encontrado nas capabilities ativas. Fazendo fallback.", decision)
			resolvedAgentID = uuid.Nil
		}
	}

	if resolvedAgentID == uuid.Nil {
		// Busca a primeira chave de API ativa para ser usada como roteador (Fallback: mock-key)
		apiKey := "sk-mock-key"
		modelName := "gpt-4o-mini"
		
		var dbKey models.APIKey
		if err := c.DB.Where("is_active = ? AND provider IN ?", true, []string{"openai", "anthropic", "openrouter", "OpenRouter"}).First(&dbKey).Error; err == nil {
			apiKey = dbKey.Key
		} else {
			log.Printf("⚠️ [Coordinator] Nenhuma chave de API ativa encontrada. Usando chave mockada.")
		}

		// Instancia LLM Rápida (Apenas para roteamento - Planner)
		routerLLM, err := llm.NewLLMProvider(modelName, apiKey)
		if err != nil {
			log.Printf("❌ [Coordinator] Erro ao instanciar Router LLM: %v", err)
			return
		}

		optionsText := ""
		for _, cap := range caps {
			optionsText += fmt.Sprintf("- '%s' (%s)\n", cap.AgentID, cap.Description)
		}

		systemPrompt := "Você é o Coordenador Central de um Swarm. Leia a mensagem do usuário e decida qual agente especialista deve assumir a tarefa. Responda apenas com o ID exato de um dos agentes abaixo:\n"
		if optionsText != "" {
			systemPrompt += optionsText
		} else {
			// Fallback dinâmico: Se não há agentes com capabilities, não há o que rotear
			log.Printf("⚠️ [Coordinator] Não há agentes no Registry para roteamento.")
			return
		}

		req := models.LLMRequest{
			Model:       "gpt-4o-mini",
			System:      systemPrompt,
			Temperature: 0.1,
			MaxTokens:   50,
			Messages: []models.LLMMessage{
				{Role: "user", Content: payload.Content},
			},
		}

		resp, err := routerLLM.Generate(context.Background(), req)
		
		if err == nil {
			decision = strings.ToLower(strings.TrimSpace(resp.Content))
		} else {
			log.Printf("⚠️ [Coordinator] Falha no Roteamento LLM: %v", err)
		}

		for _, cap := range caps {
			if cap.AgentID == decision {
				targetAgent = cap.Subject
				break
			}
		}

		if targetAgent == "" && len(caps) > 0 {
			// Fallback para o primeiro agente disponível se o LLM alucinou ou falhou
			targetAgent = caps[0].Subject
			decision = caps[0].AgentID
			resolvedAgentID, _ = uuid.Parse(decision) // Atualiza o resolvedAgentID
			
			suggested := ""
			if resp != nil {
				suggested = resp.Content
			}
			log.Printf("⚠️ [Coordinator] Fallback acionado (Sugestão anterior: '%s', Agente escolhido: '%s')", suggested, decision)
		}

		if targetAgent == "" {
			log.Printf("❌ [Coordinator] Erro fatal: Nenhum alvo encontrado para roteamento.")
			return
		}
	}

	log.Printf("🔀 [Coordinator] Roteando mensagem de '%s' para '%s' (%s)", payload.Source, decision, targetAgent)

	taskUUID := uuid.New()

	// Envia o payload encapsulado para o Especialista
	startEvent := events.AgentStartedEvent{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New(),
			EventType: events.EventAgentStarted,
			Timestamp: time.Now(),
		},
		AgentID: resolvedAgentID,
		TaskID:  taskUUID,
		Payload: string(msg.Data), // Repassa o JSON inteiro original
	}

	eventData, err := json.Marshal(startEvent)
	if err != nil {
		log.Printf("❌ [Coordinator] Erro ao serializar AgentStartedEvent: %v", err)
		return
	}
	if err := c.EventBus.Publish(targetAgent, eventData); err != nil {
		log.Printf("❌ [Coordinator] Erro ao publicar para %s: %v", targetAgent, err)
	}
}
