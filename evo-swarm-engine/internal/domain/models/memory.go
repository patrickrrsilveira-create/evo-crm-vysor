package models

import (
	"time"

	"gorm.io/gorm"
)

// ConversationMessage armazena o histórico longo (Long Term Memory) no PostgreSQL.
// Isso servirá para Analytics, Replay e interface visual do CRM, desonerando o Redis.
type ConversationMessage struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ConversationID string         `gorm:"index;not null" json:"conversation_id"` // ID externo (ex: wamidxxx)
	Role           string         `gorm:"not null" json:"role"`                  // "user", "assistant", "system", "tool"
	Content        string         `gorm:"type:text" json:"content"`
	ToolCallID     string         `json:"tool_call_id,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
