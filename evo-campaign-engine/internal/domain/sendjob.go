package domain

import (
	"time"

	"github.com/google/uuid"
)

type JobState string

const (
	JobQueued    JobState = "queued"
	JobScheduled JobState = "scheduled"
	JobSending   JobState = "sending"
	JobSent      JobState = "sent"
	JobDelivered JobState = "delivered"
	JobRead      JobState = "read"
	JobFailed    JobState = "failed"
	JobSkipped   JobState = "skipped"
)

type SendJob struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CampaignID       uuid.UUID  `gorm:"type:uuid;not null" json:"campaign_id"`
	ContactID        int        `gorm:"not null" json:"contact_id"`
	SenderInstanceID *uuid.UUID `gorm:"type:uuid" json:"sender_instance_id,omitempty"`
	Recipient        string     `gorm:"type:varchar(120);not null" json:"recipient"`
	VariantID        *uuid.UUID `gorm:"type:uuid" json:"variant_id,omitempty"`
	State            JobState   `gorm:"type:varchar(20);not null;default:'queued'" json:"state"`

	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	ReadAt      *time.Time `json:"read_at,omitempty"`

	Attempts  int    `gorm:"not null;default:0" json:"attempts"`
	LastError string `gorm:"type:text" json:"last_error,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SendJob) TableName() string { return "cmp_send_jobs" }

type DeliveryEvent struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SendJobID      uuid.UUID `gorm:"type:uuid;not null" json:"send_job_id"`
	Status         string    `gorm:"type:varchar(30);not null" json:"status"`
	ProviderStatus string    `gorm:"type:varchar(60)" json:"provider_status,omitempty"`
	Ts             time.Time `gorm:"not null;default:now()" json:"ts"`
}

func (DeliveryEvent) TableName() string { return "cmp_delivery_events" }
