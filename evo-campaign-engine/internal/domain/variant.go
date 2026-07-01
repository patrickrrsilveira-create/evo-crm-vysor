package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageVariant struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CampaignID uuid.UUID `gorm:"type:uuid;not null" json:"campaign_id"`
	Subject    string    `gorm:"type:text" json:"subject,omitempty"`
	Body       string    `gorm:"type:text;not null" json:"body"`
	MediaURL   string    `gorm:"type:text" json:"media_url,omitempty"`
	MediaType  string    `gorm:"type:varchar(30)" json:"media_type,omitempty"`
	Weight     int       `gorm:"not null;default:1" json:"weight"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (MessageVariant) TableName() string { return "cmp_message_variants" }
