package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OAuthSession armazena temporariamente o estado da sessão de login PKCE para prevenir CSRF e permitir
// a troca segura do código pelo token.
type OAuthSession struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	State     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Provider  string    `gorm:"type:varchar(50);not null"`
	AgentID   string    `gorm:"type:varchar(50);not null"`  // Agente que iniciou a requisição
	Verifier  string    `gorm:"type:varchar(255);not null"` // Code Verifier do PKCE
	CreatedAt time.Time
	ExpiresAt time.Time `gorm:"not null"` // Sessões devem expirar rápido (ex: 10 minutos)
}

// BeforeCreate garante UUID caso não venha
func (s *OAuthSession) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
