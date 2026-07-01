package domain

import (
	"time"

	"github.com/google/uuid"
)

type Suppression struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID   int       `gorm:"not null" json:"account_id"`
	ContactID   int       `gorm:"not null" json:"contact_id"`
	ChannelType string    `gorm:"type:varchar(30);not null" json:"channel_type"`
	Reason      string    `gorm:"type:varchar(60);not null;default:'opt_out'" json:"reason"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Suppression) TableName() string { return "cmp_suppression" }
