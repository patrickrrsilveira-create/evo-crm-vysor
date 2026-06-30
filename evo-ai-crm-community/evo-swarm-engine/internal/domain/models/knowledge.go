package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// KnowledgeBase representa uma base de conhecimento atrelada a uma conta.
type KnowledgeBase struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	AccountID   int64     `json:"account_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (KnowledgeBase) TableName() string {
	return "knowledge_bases"
}

// KnowledgeDocument representa um documento/site indexado na base.
type KnowledgeDocument struct {
	ID              string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	KnowledgeBaseID string    `json:"knowledge_base_id"`
	Title           string    `json:"title"`
	SourceType      string    `json:"source_type"` // ex: "pdf", "url"
	SourceURL       string    `json:"source_url"`
	Status          string    `json:"status"` // ex: "processing", "completed", "failed"
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (KnowledgeDocument) TableName() string {
	return "knowledge_documents"
}

// KnowledgeDocumentChunk armazena o texto fragmentado ("chunk") e sua respectiva assinatura vetorial.
type KnowledgeDocumentChunk struct {
	ID                  int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	KnowledgeDocumentID string `json:"knowledge_document_id"`
	Content             string `gorm:"type:text" json:"content"`
	// pgvector nativo em Go
	Embedding pgvector.Vector `gorm:"type:vector(1536)" json:"-"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (KnowledgeDocumentChunk) TableName() string {
	return "knowledge_document_chunks"
}

// KnowledgeBaseAgentBot é a tabela de junção (N:N) que atrela Bases aos Agentes
type KnowledgeBaseAgentBot struct {
	KnowledgeBaseID string    `gorm:"primaryKey" json:"knowledge_base_id"`
	AgentBotID      uuid.UUID `gorm:"primaryKey;type:uuid" json:"agent_bot_id"`
	CreatedAt       time.Time `json:"created_at"`
}

func (KnowledgeBaseAgentBot) TableName() string {
	return "knowledge_base_agent_bots"
}

// KnowledgeBaseAiAgent é a tabela de junção que atrela Bases aos Agentes de IA (EvoCore)
type KnowledgeBaseAiAgent struct {
	ID              string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	KnowledgeBaseID string    `json:"knowledge_base_id"`
	AiAgentID       uuid.UUID `gorm:"type:uuid" json:"ai_agent_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (KnowledgeBaseAiAgent) TableName() string {
	return "knowledge_base_ai_agents"
}
