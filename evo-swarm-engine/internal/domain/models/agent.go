package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Agent representa a tabela "evo_core_agents" do banco de dados (Herdado do Python/SQLAlchemy)
type Agent struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"type:varchar;not null" json:"name"`
	Role        *string        `gorm:"type:varchar" json:"role"`
	Goal        *string        `gorm:"type:text" json:"goal"`
	Description *string        `gorm:"type:text" json:"description"`
	Type        string         `gorm:"type:varchar;not null;check:type IN ('llm', 'sequential', 'parallel', 'loop', 'a2a', 'workflow', 'crew_ai', 'task', 'external')" json:"type"`
	Model       *string        `gorm:"type:varchar;default:''" json:"model"`
	APIKeyID    *uuid.UUID     `gorm:"type:uuid" json:"api_key_id"`
	Instruction *string        `gorm:"type:text" json:"instruction"`
	CardURL     *string        `gorm:"type:varchar" json:"card_url"`
	FolderID    *uuid.UUID     `gorm:"type:uuid" json:"folder_id"`
	Config      datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"config"`
	CreatedAt   time.Time      `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at"`
}

// TableName define o nome exato da tabela no PostgreSQL
func (Agent) TableName() string {
	return "evo_core_agents"
}
