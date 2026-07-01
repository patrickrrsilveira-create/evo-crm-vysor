package domain

import (
	"time"

	"github.com/google/uuid"
)

type SenderStatus string

const (
	SenderActive SenderStatus = "active"
	SenderPaused SenderStatus = "paused"
	SenderBanned SenderStatus = "banned"
)

type SenderInstance struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	InboxID      int          `gorm:"not null" json:"inbox_id"`
	ChannelType  string       `gorm:"type:varchar(30);not null" json:"channel_type"`
	Identifier   string       `gorm:"type:varchar(120);not null" json:"identifier"`
	WarmupStart  *time.Time   `gorm:"column:warmup_started_at" json:"warmup_started_at,omitempty"`
	WarmupDay    int          `gorm:"column:warmup_stage_day;not null;default:0" json:"warmup_stage_day"`
	HealthScore  int          `gorm:"not null;default:100" json:"health_score"`
	Status       SenderStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	LastSentAt   *time.Time   `json:"last_sent_at,omitempty"`
	SentToday    int          `gorm:"not null;default:0" json:"sent_today"`
	CreatedAt    time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SenderInstance) TableName() string { return "cmp_sender_instances" }
