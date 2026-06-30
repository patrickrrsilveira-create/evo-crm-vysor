package events

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// EventStore gerencia a persistência de eventos no JetStream (Event Sourcing / Auditoria)
type EventStore struct {
	EventBus *EventBus
}

func NewEventStore(bus *EventBus) (*EventStore, error) {
	js, err := bus.Conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("falha ao obter JetStream context: %v", err)
	}

	// Tenta obter o Stream. Se não existir, cria.
	_, err = js.StreamInfo("EVENT_STORE")
	if err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:        "EVENT_STORE",
			Description: "Armazena eventos do sistema para auditoria, replay e analytics",
			Subjects:    []string{"*.events.>"}, // Ex: agent.events.crm, system.events.started
			Retention:   nats.LimitsPolicy,
			MaxAge:      30 * 24 * time.Hour, // Mantém histórico por 30 dias (Ajustável)
			Storage:     nats.FileStorage,    // Salva no disco (Durabilidade de banco de dados)
		})
		if err != nil {
			return nil, fmt.Errorf("falha ao criar Stream EVENT_STORE: %v", err)
		}
		log.Println("🗄️ [EventStore] JetStream Stream 'EVENT_STORE' criado com sucesso.")
	} else {
		log.Println("🗄️ [EventStore] JetStream Stream 'EVENT_STORE' já existe.")
	}

	return &EventStore{
		EventBus: bus,
	}, nil
}

// EventLog representa o formato padrão de log de auditoria
type EventLog struct {
	Event          string      `json:"event"`
	ConversationID string      `json:"conversation_id"`
	Timestamp      string      `json:"timestamp"`
	Payload        interface{} `json:"payload"`
}

// LogEvent é uma função de helper que formata o evento e dispara no barramento (que será capturado pelo EVENT_STORE passivamente)
func (s *EventStore) LogEvent(subject string, conversationID string, eventName string, payload interface{}) {
	eventData := EventLog{
		Event:          eventName,
		ConversationID: conversationID,
		Timestamp:      time.Now().Format(time.RFC3339Nano),
		Payload:        payload,
	}

	bytes, err := json.Marshal(eventData)
	if err != nil {
		log.Printf("❌ [EventStore] Erro ao serializar evento %s: %v", eventName, err)
		return
	}
	
	// O barramento NATS normal publica. Como o tópico "subject" bate com "*.events.>", o JetStream grava no HD.
	if err := s.EventBus.Publish(subject, bytes); err != nil {
		log.Printf("❌ [EventStore] Erro ao publicar evento %s no NATS: %v", eventName, err)
	}
}
