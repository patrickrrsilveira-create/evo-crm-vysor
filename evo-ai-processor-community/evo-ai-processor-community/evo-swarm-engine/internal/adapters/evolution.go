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

		// Validação básica se é evento de mensagem nova via WhatsApp (Suporta Evolution v1 e v2)
		if eventStr == "messages.upsert" || eventStr == "Message" {
			log.Println("📥 [EvolutionWebhook] Nova mensagem do WhatsApp recebida! Disparando para o Swarm...")

			// Extração segura de Sender e Content
			sender := ""
			content := ""
			
			if data, ok := payload["data"].(map[string]interface{}); ok {
				// Formato v2 ("Message")
				if info, ok := data["Info"].(map[string]interface{}); ok {
					if s, ok := info["Sender"].(string); ok {
						sender = s
					}
				}
				if msg, ok := data["Message"].(map[string]interface{}); ok {
					if c, ok := msg["conversation"].(string); ok {
						content = c
					} else if extMsg, ok := msg["extendedTextMessage"].(map[string]interface{}); ok {
						if c, ok := extMsg["text"].(string); ok {
							content = c
						}
					}
				}
				
				// Formato v1 ("messages.upsert")
				if msgs, ok := data["messages"].([]interface{}); ok && len(msgs) > 0 {
					if firstMsg, ok := msgs[0].(map[string]interface{}); ok {
						if pushName, ok := firstMsg["pushName"].(string); ok && sender == "" {
							sender = pushName // Fallback
						}
						if key, ok := firstMsg["key"].(map[string]interface{}); ok && sender == "" {
							if s, ok := key["remoteJid"].(string); ok {
								sender = s
							}
						}
						if msg, ok := firstMsg["message"].(map[string]interface{}); ok {
							if c, ok := msg["conversation"].(string); ok {
								content = c
							} else if extMsg, ok := msg["extendedTextMessage"].(map[string]interface{}); ok {
								if c, ok := extMsg["text"].(string); ok {
									content = c
								}
							}
						}
					}
				}
			}

			// Se não conseguiu extrair, usa raw JSON para debug
			if content == "" {
				rawBytes, _ := json.Marshal(payload)
				content = "[Raw Payload] " + string(rawBytes)
			}
			if sender == "" {
				sender = "unknown_whatsapp_sender"
			}

			eventData, _ := json.Marshal(map[string]interface{}{
				"source":   "whatsapp_evolution",
				"content":  content,
				"sender":   sender,
				"agent_id": "", // Evolution nativo não sabe o AgentID, Coordinator fará fallback
				"payload":  payload,
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
