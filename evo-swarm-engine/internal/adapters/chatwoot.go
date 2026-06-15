package adapters

import (
	"context"
	"encoding/json"
	"log"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/database"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/models"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
)

// ChatwootMirrorAdapter escuta todas as mensagens enviadas pelos agentes
// e faz o "Mirroring" (espelhamento) direto no banco PostgreSQL do Chatwoot.
type ChatwootMirrorAdapter struct {
	EventBus *evbus.EventBus
}

func NewChatwootMirrorAdapter(bus *evbus.EventBus) *ChatwootMirrorAdapter {
	return &ChatwootMirrorAdapter{
		EventBus: bus,
	}
}

func (a *ChatwootMirrorAdapter) Name() string {
	return "chatwoot_db_mirror"
}

// Start inicia a inscrição no barramento de eventos (NATS)
func (a *ChatwootMirrorAdapter) Start(ctx context.Context) error {
	log.Println("🪞 [ChatwootMirror] Inicializado - Escutando eventos 'message.sent' para espelhamento nativo (PostgreSQL)")

	_, err := a.EventBus.Subscribe(string(events.EventMessageSent), a.handleMessageSent)
	return err
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

	result := database.DB.Create(&chatwootMsg)
	if result.Error != nil {
		log.Printf("❌ [ChatwootMirror] Falha ao espelhar mensagem no DB: %v", result.Error)
		return
	}

	log.Printf("🪞 [ChatwootMirror] Mensagem inserida com sucesso! (ID: %d)", chatwootMsg.ID)
}

// RegisterWebhookRoute registra a rota HTTP para receber eventos do Chatwoot (Incoming)
func (a *ChatwootMirrorAdapter) RegisterWebhookRoute(app *fiber.App) {
	app.Post("/webhooks/chatwoot", func(c *fiber.Ctx) error {
		// Chatwoot envia o payload quando uma mensagem é criada
		var payload map[string]interface{}
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).SendString("Bad Request")
		}

		// Checa se é um evento de mensagem criada e se não foi enviada por nós
		if eventType, ok := payload["event"].(string); ok && eventType == "message_created" {
			if msgType, ok := payload["message_type"].(float64); ok && msgType == 0 { // 0 = incoming (user)
				log.Println("📥 [ChatwootWebhook] Nova mensagem recebida do usuário! Disparando para o Swarm...")

				// Dispara no NATS indicando Input do Usuário
				// Em produção real usaríamos uma struct tipada de payload
				eventData, _ := json.Marshal(map[string]interface{}{
					"source":  "chatwoot",
					"content": payload["content"],
				})

				a.EventBus.Publish(string(events.EventMessageReceived), eventData)
			}
		}

		return c.SendStatus(200)
	})
}
