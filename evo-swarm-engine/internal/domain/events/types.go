package events

import (
	"time"

	"github.com/google/uuid"
)

// EventType define os tipos estritos de eventos que circulam no barramento
type EventType string

const (
	EventAgentStarted    EventType = "agent.started"
	EventAgentFinished   EventType = "agent.finished"
	EventMessageSent     EventType = "message.sent"     // Saída (Bot -> Usuário)
	EventMessageReceived EventType = "message.received" // Entrada (Usuário -> Bot)
	EventLeadCreated     EventType = "lead.created"
)

// BaseEvent é a estrutura comum a todos os eventos assíncronos
type BaseEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	EventType EventType `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	TraceID   string    `json:"trace_id,omitempty"` // Para observabilidade distribuída
}

// AgentStartedEvent é disparado quando um Agente Especialista inicia uma task
type AgentStartedEvent struct {
	BaseEvent
	AgentID   uuid.UUID `json:"agent_id"`
	AgentName string    `json:"agent_name"`
	TaskID    uuid.UUID `json:"task_id"`
	Payload   string    `json:"payload"`
}

// AgentFinishedEvent é disparado quando o Agente conclui com sucesso ou falha
type AgentFinishedEvent struct {
	BaseEvent
	AgentID uuid.UUID `json:"agent_id"`
	TaskID  uuid.UUID `json:"task_id"`
	Result  string    `json:"result"`
	Success bool      `json:"success"`
}
