package models

import (
	"github.com/google/uuid"
	"time"
)

// Estruturas Oficiais do Protocolo A2A (Google)

// A2ARequest representa uma chamada JSON-RPC padrão
type A2ARequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      string                 `json:"id"`
}

// A2AResponse representa o retorno JSON-RPC
type A2AResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *A2AError   `json:"error,omitempty"`
	ID      string      `json:"id"`
}

type A2AError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// A2ATask representa um objeto de Resposta de Tarefa (Message Send)
type A2ATask struct {
	ID        string        `json:"id"`
	ContextID string        `json:"contextId"`
	Status    A2ATaskStatus `json:"status"`
	Artifacts []A2AArtifact `json:"artifacts"`
	Kind      string        `json:"kind"` // Sempre "task"
	History   []A2AMessage  `json:"history,omitempty"`
}

type A2ATaskStatus struct {
	State     string    `json:"state"` // "completed", "failed", "pending"
	Timestamp time.Time `json:"timestamp"`
}

type A2AArtifact struct {
	ArtifactID string    `json:"artifactId"`
	Parts      []A2APart `json:"parts"`
}

type A2APart struct {
	Type string `json:"type"` // "text", "file"
	Text string `json:"text,omitempty"`
}

type A2AMessage struct {
	Role      string    `json:"role"`
	Parts     []A2APart `json:"parts"`
	MessageID uuid.UUID `json:"messageId,omitempty"`
	TaskID    string    `json:"taskId,omitempty"`
	ContextID string    `json:"contextId,omitempty"`
	Kind      string    `json:"kind"` // "message"
	Timestamp time.Time `json:"timestamp,omitempty"`
}
