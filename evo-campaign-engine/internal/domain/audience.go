package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AudienceMode string

const (
	AudienceStatic  AudienceMode = "static"
	AudienceDynamic AudienceMode = "dynamic"
)

type CampaignAudience struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CampaignID    uuid.UUID      `gorm:"type:uuid;not null" json:"campaign_id"`
	Mode          AudienceMode   `gorm:"type:varchar(20);not null;default:'static'" json:"mode"`
	SegmentFilter datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"segment_filter"`
	TotalContacts int            `gorm:"not null;default:0" json:"total_contacts"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

func (CampaignAudience) TableName() string { return "cmp_campaign_audience" }

type MemberState string

const (
	MemberPending MemberState = "pending"
	MemberQueued  MemberState = "queued"
	MemberDone    MemberState = "done"
	MemberSkipped MemberState = "skipped"
)

type AudienceMember struct {
	ID         uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CampaignID uuid.UUID   `gorm:"type:uuid;not null" json:"campaign_id"`
	ContactID  int         `gorm:"not null" json:"contact_id"`
	Recipient  string      `gorm:"type:varchar(120);not null" json:"recipient"`
	Timezone   string      `gorm:"type:varchar(60)" json:"timezone,omitempty"`
	State      MemberState `gorm:"type:varchar(20);not null;default:'pending'" json:"state"`
}

func (AudienceMember) TableName() string { return "cmp_audience_members" }
