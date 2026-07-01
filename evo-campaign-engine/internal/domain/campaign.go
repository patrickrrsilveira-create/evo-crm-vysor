package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CampaignStatus string

const (
	CampaignDraft     CampaignStatus = "draft"
	CampaignScheduled CampaignStatus = "scheduled"
	CampaignRunning   CampaignStatus = "running"
	CampaignPaused    CampaignStatus = "paused"
	CampaignCompleted CampaignStatus = "completed"
	CampaignCancelled CampaignStatus = "cancelled"
)

type Campaign struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID         int            `gorm:"not null" json:"account_id"`
	Name              string         `gorm:"type:varchar(255);not null" json:"name"`
	TriggerType       string         `gorm:"type:varchar(20);not null;default:'manual'" json:"trigger_type"`
	TriggerConfig     datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"trigger_config"`
	ThrottleProfileID *uuid.UUID     `gorm:"type:uuid" json:"throttle_profile_id,omitempty"`
	Status            CampaignStatus `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`

	TotalRecipients int `gorm:"not null;default:0" json:"total_recipients"`
	SentCount       int `gorm:"not null;default:0" json:"sent_count"`
	DeliveredCount  int `gorm:"not null;default:0" json:"delivered_count"`
	FailedCount     int `gorm:"not null;default:0" json:"failed_count"`

	CreatedBy  *int       `json:"created_by,omitempty"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`

	Channels []CampaignChannel `gorm:"foreignKey:CampaignID" json:"channels,omitempty"`
	Variants []MessageVariant  `gorm:"foreignKey:CampaignID" json:"variants,omitempty"`
	Audience *CampaignAudience `gorm:"foreignKey:CampaignID" json:"audience,omitempty"`
}

func (Campaign) TableName() string { return "cmp_campaigns" }

type CampaignChannel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CampaignID uuid.UUID `gorm:"type:uuid;not null" json:"campaign_id"`
	InboxID    int       `gorm:"not null" json:"inbox_id"`
}

func (CampaignChannel) TableName() string { return "cmp_campaign_channels" }
