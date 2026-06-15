package adapters

import (
	"context"
	"log"
)

type CalendarAdapter struct {
	ServiceAccountJSON string
}

func NewCalendarAdapter(serviceAccount string) *CalendarAdapter {
	return &CalendarAdapter{
		ServiceAccountJSON: serviceAccount,
	}
}

func (a *CalendarAdapter) Name() string {
	return "google_calendar"
}

// Escuta via Webhooks (Push notifications) de atualizações na agenda
func (a *CalendarAdapter) Start(ctx context.Context) error {
	log.Println("📅 [CalendarAdapter] Inicializado - Escutando eventos do Google Calendar")
	return nil
}

func (a *CalendarAdapter) SendMessage(ctx context.Context, to string, content string) error {
	log.Printf("📅 [CalendarAdapter] Operação no calendário (convite/notificação) para: %s", to)
	return nil
}
