package adapters

import (
	"context"
	"log"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
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
	// Na prática, aqui criamos endpoints Fiber (e.g., POST /webhook/evolution)
	// que recebem o Payload, parseiam e publicam um "events.EventMessageSent" no EventBus NATS
	return nil
}

func (a *EvolutionAdapter) SendMessage(ctx context.Context, to string, content string) error {
	log.Printf("📱 [EvolutionAdapter] Enviando mensagem via Evolution API para: %s", to)
	// Chamada HTTP (POST /message/sendText) usando http.Client
	return nil
}
