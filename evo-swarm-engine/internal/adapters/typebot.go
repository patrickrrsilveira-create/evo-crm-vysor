package adapters

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
)

// TypebotAdapter funciona como uma ponte síncrona (HTTP Block)
// conectando as interfaces do Typebot ao motor assíncrono do Swarm.
type TypebotAdapter struct {
	EventBus *evbus.EventBus
}

func NewTypebotAdapter(bus *evbus.EventBus) *TypebotAdapter {
	return &TypebotAdapter{
		EventBus: bus,
	}
}

// Start não tem assinatura passiva de longa duração, pois atua sob demanda HTTP
func (a *TypebotAdapter) Start(ctx context.Context) error {
	log.Println("🤖 [TypebotAdapter] Inicializado (Ponte Síncrona via HTTP)")
	return nil
}

// RegisterWebhookRoute expõe o endpoint que o Typebot consumirá no seu bloco de "HTTP Request"
func (a *TypebotAdapter) RegisterWebhookRoute(app *fiber.App) {
	app.Post("/api/webhooks/typebot", func(c *fiber.Ctx) error {
		var payload struct {
			SessionID string `json:"sessionId"`
			Message   string `json:"message"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Payload inválido. Esperado 'sessionId' e 'message'",
			})
		}

		if payload.SessionID == "" || payload.Message == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Campos sessionId e message são obrigatórios",
			})
		}

		log.Printf("📥 [TypebotAdapter] Mensagem recebida da sessão: %s", payload.SessionID)

		// 1. Prepara um canal para receber a resposta do Swarm assincronamente
		responseChan := make(chan string)

		// Assinatura efêmera para interceptar o que o agente enviará de volta
		// Em vez do TaskID, ouvimos o outbound.message global e filtramos pelo sender
		sub, err := a.EventBus.Conn.Subscribe("outbound.message", func(msg *nats.Msg) {
			var out struct {
				Sender  string `json:"sender"`
				Content string `json:"content"`
			}
			if err := json.Unmarshal(msg.Data, &out); err == nil {
				if out.Sender == payload.SessionID {
					responseChan <- out.Content
				}
			}
		})
		if err != nil {
			log.Printf("❌ [TypebotAdapter] Erro ao assinar stream de resposta: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Error")
		}
		defer sub.Unsubscribe() // Limpa a memória assim que a função terminar

		// 2. Dispara a mensagem para o Coordinator rotear
		eventData, _ := json.Marshal(map[string]interface{}{
			"source":     "typebot",
			"sender":     payload.SessionID, // O session ID atua como sender único
			"context_id": payload.SessionID,
			"content":    payload.Message,
		})

		if err := a.EventBus.Publish(string(events.EventMessageReceived), eventData); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Erro ao enviar para o motor")
		}

		// 3. Aguarda a resposta (Com Timeout rigoroso por causa do Typebot)
		select {
		case replyContent := <-responseChan:
			log.Printf("📤 [TypebotAdapter] Resposta devolvida com sucesso para sessão %s", payload.SessionID)
			return c.JSON(fiber.Map{
				"reply": replyContent,
			})
		case <-time.After(25 * time.Second): // O Typebot costuma cancelar após 30s
			log.Printf("⏳ [TypebotAdapter] Timeout aguardando IA para a sessão %s", payload.SessionID)
			return c.JSON(fiber.Map{
				"reply": "Desculpe, estou processando muita informação. Pode aguardar um momento e dizer 'continuar'?",
			})
		}
	})
}
