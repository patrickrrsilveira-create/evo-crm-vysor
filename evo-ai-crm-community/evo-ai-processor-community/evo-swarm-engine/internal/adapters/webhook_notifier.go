package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/domain/events"
	evbus "github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
	"github.com/nats-io/nats.go"
)

// WebhookNotifierAdapter ouve eventos de sistema (como handoff) e notifica URLs externas
type WebhookNotifierAdapter struct {
	EventBus   *evbus.EventBus
	HTTPClient *http.Client
}

func NewWebhookNotifierAdapter(bus *evbus.EventBus) *WebhookNotifierAdapter {
	return &WebhookNotifierAdapter{
		EventBus: bus,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *WebhookNotifierAdapter) Start(ctx context.Context) error {
	log.Println("🌐 [WebhookNotifier] Iniciando ouvinte de webhooks (Handoff)...")

	_, err := a.EventBus.Conn.QueueSubscribe(string(events.EventAgentHandoff), "webhook_notifier_group", a.handleHandoff)
	return err
}

func (a *WebhookNotifierAdapter) handleHandoff(msg *nats.Msg) {
	webhookURL := os.Getenv("WEBHOOK_HANDOFF_URL")
	if webhookURL == "" {
		log.Println("⚠️ [WebhookNotifier] Evento Handoff recebido, mas WEBHOOK_HANDOFF_URL não está configurada.")
		return
	}

	var handoffEvent events.AgentHandoffEvent
	if err := json.Unmarshal(msg.Data, &handoffEvent); err != nil {
		log.Printf("❌ [WebhookNotifier] Erro ao decodificar AgentHandoffEvent: %v", err)
		return
	}

	log.Printf("🚀 [WebhookNotifier] Disparando Webhook POST para %s", webhookURL)

	payloadBytes, err := json.Marshal(handoffEvent)
	if err != nil {
		log.Printf("❌ [WebhookNotifier] Erro ao codificar JSON para webhook: %v", err)
		return
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("❌ [WebhookNotifier] Erro ao criar request HTTP: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Evo-Swarm-Engine/1.0")

	// Disparo não bloqueante
	go func() {
		resp, err := a.HTTPClient.Do(req)
		if err != nil {
			log.Printf("❌ [WebhookNotifier] Falha ao enviar webhook para %s: %v", webhookURL, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("✅ [WebhookNotifier] Webhook entregue com sucesso (Status %d)", resp.StatusCode)
		} else {
			log.Printf("⚠️ [WebhookNotifier] Webhook retornou status não-sucesso (Status %d)", resp.StatusCode)
		}
	}()
}
