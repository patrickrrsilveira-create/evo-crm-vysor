package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AgentIntegration mapeia a tabela evo_core_agent_integrations legada do Python
type AgentIntegration struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AgentID   uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_agent_provider" json:"agent_id"`
	Provider  string         `gorm:"type:varchar;not null;uniqueIndex:idx_agent_provider" json:"provider"`
	Config    datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"config"`
	IsActive  bool           `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updated_at"`
}

// TableName garante a compatibilidade com a tabela legada do Python/SQLAlchemy
func (AgentIntegration) TableName() string {
	return "evo_core_agent_integrations"
}
