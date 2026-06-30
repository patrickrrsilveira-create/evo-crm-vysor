package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

// ChatwootMirrorAdapter escuta todas as mensagens enviadas pelos agentes
// e faz o "Mirroring" (espelhamento) direto no banco PostgreSQL do Chatwoot.
// Também escuta outbound.message para enviar respostas via REST API do CRM.
type ChatwootMirrorAdapter struct {
	EventBus *evbus.EventBus
	db       *gorm.DB
}

func NewChatwootMirrorAdapter(bus *evbus.EventBus, db *gorm.DB) *ChatwootMirrorAdapter {
	return &ChatwootMirrorAdapter{
		EventBus: bus,
		db:       db,
	}
}

func (a *ChatwootMirrorAdapter) Name() string {
	return "chatwoot_db_mirror"
}

// Start inicia a inscrição no barramento de eventos (NATS)
func (a *ChatwootMirrorAdapter) Start(ctx context.Context) error {
	log.Println("🪞 [ChatwootMirror] Inicializado - Escutando eventos 'message.sent' para espelhamento nativo (PostgreSQL)")

	_, err := a.EventBus.Subscribe(string(events.EventMessageSent), a.handleMessageSent)
	if err != nil {
		return err
	}

	// Escuta outbound.message para enviar respostas de volta ao Chatwoot via REST API
	_, err = a.EventBus.Conn.QueueSubscribe("outbound.message", "chatwoot_outbound_group", a.handleOutboundMessage)
	if err != nil {
		log.Printf("⚠️ [ChatwootMirror] Falha ao assinar outbound.message: %v", err)
	} else {
		log.Println("📤 [ChatwootMirror] Escutando outbound.message para envio via REST API do CRM")
	}

	return nil
}

// SendMessage não tem uso direto aqui, pois este adaptador é estritamente para espelhamento passivo
func (a *ChatwootMirrorAdapter) SendMessage(ctx context.Context, to string, content string) error {
	return nil
}

// handleMessageSent intercepta o evento de mensagem e injeta no Chatwoot
func (a *ChatwootMirrorAdapter) handleMessageSent(msg *nats.Msg) {
	// A estrutura abaixo seria definida no types.go real, aqui simulamos uma genérica
	var payload struct {
		ConversationID int64  `json:"conversation_id"`
		AccountID      int64  `json:"account_id"`
		InboxID        int64  `json:"inbox_id"`
		Content        string `json:"content"`
	}

	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		log.Printf("[ChatwootMirror] Erro ao decodificar evento: %v", err)
		return
	}

	// Insere diretamente na tabela `messages` do Chatwoot usando GORM
	chatwootMsg := models.ChatwootMessage{
		AccountID:      payload.AccountID,
		InboxID:        payload.InboxID,
		ConversationID: payload.ConversationID,
		Content:        payload.Content,
		MessageType:    1,          // 1 = outgoing (Enviado pelo nosso Agente/Bot)
		SenderType:     "AgentBot", // Diferencia a UI do Chatwoot
		SenderID:       1,          // ID do bot na tabela users/agent_bots
		Status:         1,          // 1 = delivered
	}

	result := a.db.Create(&chatwootMsg)
	if result.Error != nil {
		log.Printf("❌ [ChatwootMirror] Falha ao espelhar mensagem no DB: %v", result.Error)
		return
	}

	log.Printf("🪞 [ChatwootMirror] Mensagem inserida com sucesso! (ID: %d)", chatwootMsg.ID)
}

// handleOutboundMessage recebe a resposta do agente e envia de volta ao Chatwoot via REST API
func (a *ChatwootMirrorAdapter) handleOutboundMessage(msg *nats.Msg) {
	var outbound struct {
		Source         string `json:"source"`
		ConversationID string `json:"conversation_id"`
		AccountID      int64  `json:"account_id"`
		Content        string `json:"content"`
	}

	if err := json.Unmarshal(msg.Data, &outbound); err != nil {
		log.Printf("❌ [ChatwootOutbound] Erro ao decodificar outbound.message: %v", err)
		return
	}

	// Só processa mensagens que vieram originalmente do Chatwoot
	if outbound.Source != "chatwoot" {
		return
	}

	log.Printf("📤 [ChatwootOutbound] Enviando resposta do agente de volta ao CRM (ConversationID: %s)", outbound.ConversationID)

	// Envia via REST API do CRM (Chatwoot interno)
	crmURL := os.Getenv("EVO_AI_CRM_URL")
	if crmURL == "" {
		crmURL = "http://evo_crm:3000"
	}
	apiToken := os.Getenv("EVOAI_CRM_API_TOKEN")

	accountID := outbound.AccountID
	if accountID == 0 {
		accountID = 1 // Fallback para account padrão
	}

	// In the Evo-CRM fork, the API uses the UUID directly
	url := fmt.Sprintf("%s/api/v1/conversations/%s/messages", crmURL, outbound.ConversationID)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"content":      outbound.Content,
		"message_type": "outgoing",
		"private":      false,
	})

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("❌ [ChatwootOutbound] Erro ao criar request HTTP: %v", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("api_access_token", apiToken)
	httpReq.Header.Set("X-Service-Token", apiToken)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("❌ [ChatwootOutbound] Erro ao enviar mensagem para CRM: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("✅ [ChatwootOutbound] Mensagem enviada com sucesso ao CRM! (Status: %d)", resp.StatusCode)
	} else {
		log.Printf("❌ [ChatwootOutbound] CRM respondeu com erro (Status: %d): %s", resp.StatusCode, string(body))
	}
}

// RegisterWebhookRoute registra a rota HTTP para receber eventos do Chatwoot (Incoming)
func (a *ChatwootMirrorAdapter) RegisterWebhookRoute(app *fiber.App) {
	app.Post("/webhooks/chatwoot", func(c *fiber.Ctx) error {
		// Chatwoot envia o payload quando uma mensagem é criada
		var payload map[string]interface{}
		if err := c.BodyParser(&payload); err != nil {
			log.Printf("❌ [ChatwootWebhook] Erro ao parsear body: %v", err)
			return c.Status(400).SendString("Bad Request")
		}

		// Debug: log dos campos-chave do payload
		eventType, _ := payload["event"].(string)
		log.Printf("📥 [ChatwootWebhook] Evento recebido: event='%s', message_type='%v' (tipo: %T)", eventType, payload["message_type"], payload["message_type"])

		// Checa se é um evento de mensagem criada
		if eventType == "message_created" {
			// message_type pode vir como float64 (0) OU como string ("incoming")
			isIncoming := false
			switch mt := payload["message_type"].(type) {
			case float64:
				isIncoming = mt == 0
			case string:
				isIncoming = strings.EqualFold(mt, "incoming")
			}

			if isIncoming {
				log.Println("📥 [ChatwootWebhook] Nova mensagem recebida do usuário! Disparando para o Swarm...")

				// Extrai o agent_id da Query String da URL configurada no AgentBot/Webhook do Chatwoot
				agentID := c.Query("agent_id", "")

				// Extrai o conteúdo da mensagem (pode estar em payload.content ou payload.content.text)
				content := ""
				if c, ok := payload["content"].(string); ok {
					content = c
				}

				// Extrai conversation_id e account_id para a resposta de volta
				var conversationID string
				var accountID int64
				if conv, ok := payload["conversation"].(map[string]interface{}); ok {
					if idStr, ok := conv["id"].(string); ok {
						conversationID = idStr
					}
					if accID, ok := conv["account_id"].(float64); ok {
						accountID = int64(accID)
					}
				}

				// Extrai sender para identificação
				sender := ""
				if senderObj, ok := payload["sender"].(map[string]interface{}); ok {
					if name, ok := senderObj["name"].(string); ok {
						sender = name
					}
					if id, ok := senderObj["id"].(float64); ok && sender == "" {
						sender = fmt.Sprintf("user_%d", int64(id))
					}
				}

				log.Printf("📥 [ChatwootWebhook] Detalhes: content='%s', agent_id='%s', conversation_id='%s', account_id=%d, sender='%s'",
					content, agentID, conversationID, accountID, sender)

				// Dispara no NATS indicando Input do Usuário
				eventData, _ := json.Marshal(map[string]interface{}{
					"source":          "chatwoot",
					"content":         content,
					"agent_id":        agentID,
					"sender":          sender,
					"conversation_id": conversationID,
					"account_id":      accountID,
					"payload":         payload,
				})

				a.EventBus.Publish(string(events.EventMessageReceived), eventData)
			} else {
				log.Printf("⏭️ [ChatwootWebhook] Mensagem ignorada (message_type não é incoming: %v)", payload["message_type"])
			}
		}

		return c.SendStatus(200)
	})
}
