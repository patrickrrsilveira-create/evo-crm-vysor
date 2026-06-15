package adapters

import (
	"context"
	"log"

	"github.com/PatrickRSilveira/evo-swarm-engine/internal/events"
)

// TeamsAdapter implementa ChannelAdapter para Microsoft Teams (Bot Framework / Graph API)
type TeamsAdapter struct {
	EventBus *events.EventBus
	BotAppID string
	BotPass  string
}

func NewTeamsAdapter(bus *events.EventBus, appID, pass string) *TeamsAdapter {
	return &TeamsAdapter{
		EventBus: bus,
		BotAppID: appID,
		BotPass:  pass,
	}
}

func (a *TeamsAdapter) Name() string {
	return "microsoft_teams"
}

func (a *TeamsAdapter) Start(ctx context.Context) error {
	log.Println("📞 [TeamsAdapter] Inicializado - Escutando eventos do MS Teams Bot Framework")
	// Integrar com webhook do bot framework (/api/messages)
	return nil
}

func (a *TeamsAdapter) SendMessage(ctx context.Context, to string, content string) error {
	log.Printf("📞 [TeamsAdapter] Enviando mensagem via MS Teams Bot para canal/usuário: %s", to)
	// Chamada para Bot Framework API para injetar Activity (Message)
	return nil
}
