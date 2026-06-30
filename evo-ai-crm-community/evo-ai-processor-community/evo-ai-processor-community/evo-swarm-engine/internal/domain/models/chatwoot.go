package models

import (
	"time"

	"gorm.io/gorm"
)

// ChatwootMessage representa a tabela `messages` do banco nativo do Chatwoot (PostgreSQL).
// Isso nos permite injetar mensagens dos nossos Agentes diretamente na UI do Chatwoot
// sem passar pela latência de HTTP APIs.
type ChatwootMessage struct {
	ID             int64          `gorm:"primaryKey;column:id"`
	Content        string         `gorm:"column:content;type:text"`
	AccountID      int64          `gorm:"column:account_id;index"`
	InboxID        int64          `gorm:"column:inbox_id;index"`
	ConversationID int64          `gorm:"column:conversation_id;index"`
	MessageType    int            `gorm:"column:message_type"` // 0=incoming, 1=outgoing, 2=activity, 3=template
	ContentType    string         `gorm:"column:content_type;default:'text'"`
	Status         int            `gorm:"column:status;default:0"` // 0=sent, 1=delivered, 2=read, 3=failed
	SenderType     string         `gorm:"column:sender_type"`      // 'User', 'Contact', 'AgentBot'
	SenderID       int64          `gorm:"column:sender_id"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// TableName garante que o GORM mapeie para a tabela exata do Chatwoot
func (ChatwootMessage) TableName() string {
	return "messages"
}

// ChatwootConversation representa a tabela `conversations` do Chatwoot
type ChatwootConversation struct {
	ID        int64          `gorm:"primaryKey;column:id"`
	AccountID int64          `gorm:"column:account_id;index"`
	InboxID   int64          `gorm:"column:inbox_id;index"`
	Status    int            `gorm:"column:status;default:0"` // 0=open, 1=resolved, 2=pending, 3=snoozed
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ChatwootConversation) TableName() string {
	return "conversations"
}
