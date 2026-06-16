package models

import (
	"time"

	"github.com/google/uuid"
)

// APIKey representa a tabela "evo_core_api_keys"
type APIKey struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"type:varchar;not null" json:"name"`
	Key       string    `gorm:"type:varchar;not null" json:"key"` // Criptografado ou plano, depenendo da config
	Provider  string    `gorm:"type:varchar" json:"provider"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (APIKey) TableName() string {
	return "evo_core_api_keys"
}
