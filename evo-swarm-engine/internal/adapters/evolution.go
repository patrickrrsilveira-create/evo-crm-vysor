package adapters

import (
	"context"
	"encoding/json"
	"log"

	domainEvents "github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	"github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/gofiber/fiber/v2"
)

// EvolutionAdapter implementa ChannelAdapter para a Evolution API (WhatsApp)
type EvolutionAdapter struct {
	EventBus *events.EventBus
	BaseURL  string
	APIKey   string
}

func NewEvolutionAdapter(bus *events.EventBus, url, apiKey string) *EvolutionAdapter {
	return &EvolutionAdapter{
		EventBus: bus,
		BaseURL:  url,
		APIKey:   apiKey,
	}
}

func (a *EvolutionAdapter) Name() string {
	return "whatsapp_evolution"
}

func (a *EvolutionAdapter) Start(ctx context.Context) error {
	log.Println("📱 [EvolutionAdapter] Inicializado - Webhook pronto para receber mensagens do WhatsApp")
	return nil
}

func (a *EvolutionAdapter) SendMessage(ctx context.Context, to string, content string) error {
	log.Printf("📱 [EvolutionAdapter] Enviando mensagem via Evolution API para: %s", to)
	// Chamada HTTP (POST /message/sendText) usando http.Client
	return nil
}

// RegisterWebhookRoute registra a rota HTTP para receber eventos da Evolution API (WhatsApp)
func (a *EvolutionAdapter) RegisterWebhookRoute(app *fiber.App) {
	handler := func(c *fiber.Ctx) error {
		// Blindagem de Segurança (Webhook Verification)
		// Verifica o header de autenticação configurado na Evolution API
		webhookKey := c.Get("apikey")
		if webhookKey != a.APIKey && a.APIKey != "" {
			log.Printf("⚠️ [EvolutionWebhook] Bloqueado: Tentativa de webhook com chave inválida (IP: %s)", c.IP())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Forbidden: Invalid Webhook API Key",
			})
		}

		var payload map[string]interface{}
		if err := c.BodyParser(&payload); err != nil {
			log.Printf("⚠️ [EvolutionWebhook] Erro no BodyParser: %v", err)
			return c.Status(400).SendString("Bad Request")
		}

		eventStr, _ := payload["event"].(string)
		log.Printf("📥 [EvolutionWebhook] Evento Recebido: '%s' (payload completo: %+v)", eventStr, payload)

		// Validação básica se é evento de mensagem nova via WhatsApp
		if eventStr == "messages.upsert" {
			log.Println("📥 [EvolutionWebhook] Nova mensagem do WhatsApp recebida! Disparando para o Swarm...")

			eventData, _ := json.Marshal(map[string]interface{}{
				"source":  "whatsapp_evolution",
				"payload": payload,
			})

			// Dispara evento indicando Input do Usuário
			a.EventBus.Publish(string(domainEvents.EventMessageReceived), eventData)
		}

		return c.SendStatus(200)
	}

	app.Post("/webhooks/evolution", handler)
	app.Post("/webhooks/whatsapp/evolution_go", handler)
	app.Post("/webhooks/whatsapp/evolution", handler)
}
