package adapters

import (
	"context"
	"encoding/json"
	"log"

	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/gofiber/fiber/v2"
)

// N8nAdapter atua como receptor (Inbound) de webhooks acionados pelo n8n
// Ele permite que automações do n8n disparem fluxos proativos no Swarm (ex: "Novo lead criado, analise!")
type N8nAdapter struct {
	EventBus *evbus.EventBus
}

func NewN8nAdapter(bus *evbus.EventBus) *N8nAdapter {
	return &N8nAdapter{
		EventBus: bus,
	}
}

// Start não mantém estado porque este adaptador é acionado estritamente via Webhook HTTP
func (a *N8nAdapter) Start(ctx context.Context) error {
	log.Println("⚙️ [N8nAdapter] Inicializado (Aguardando Webhooks do n8n)")
	return nil
}

// RegisterWebhookRoute expõe o endpoint POST que o nó HTTP do n8n deve chamar
func (a *N8nAdapter) RegisterWebhookRoute(app *fiber.App) {
	app.Post("/api/webhooks/n8n", func(c *fiber.Ctx) error {
		// O n8n é super dinâmico, então usamos um map solto
		var payload map[string]interface{}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Payload JSON inválido",
			})
		}

		// A lógica do nosso adaptador é: 
		// O n8n deve mandar no mínimo "event" e "payload"
		eventName, ok := payload["event"].(string)
		if !ok || eventName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Campo 'event' é obrigatório (ex: message.received, lead.created)",
			})
		}

		log.Printf("📥 [N8nAdapter] Evento recebido do n8n: %s", eventName)

		// Dispara cegamente o payload embalado para o barramento assíncrono (Coordinator se vira)
		// Isso permite máxima flexibilidade no nó do n8n
		eventData, err := json.Marshal(payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Erro interno ao empacotar evento")
		}

		// Publica no tópico mapeado (se o n8n mandou message.received, vai simular que o usuário digitou algo)
		a.EventBus.Publish(eventName, eventData)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "Event Accepted",
		})
	})
}
