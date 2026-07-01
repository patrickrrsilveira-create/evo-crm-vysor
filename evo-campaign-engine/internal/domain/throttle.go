package domain

import (
	"time"

	"github.com/google/uuid"
)

type ThrottleProfile struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(120);not null" json:"name"`
	ChannelType string    `gorm:"type:varchar(30);not null" json:"channel_type"`

	DailyCapStart    int     `gorm:"not null;default:20" json:"daily_cap_start"`
	DailyCapMax      int     `gorm:"not null;default:500" json:"daily_cap_max"`
	WarmupMultiplier float64 `gorm:"type:numeric;not null;default:1.3" json:"warmup_multiplier"`
	WarmupStepDays   int     `gorm:"not null;default:1" json:"warmup_step_days"`
	HourlyCap        int     `gorm:"not null;default:60" json:"hourly_cap"`

	MinDelaySec        int `gorm:"not null;default:45" json:"min_delay_sec"`
	MaxDelaySec        int `gorm:"not null;default:120" json:"max_delay_sec"`
	CoffeeBreakEveryN  int `gorm:"not null;default:25" json:"coffee_break_every_n"`
	CoffeeBreakMinSec  int `gorm:"not null;default:120" json:"coffee_break_min_sec"`
	CoffeeBreakMaxSec  int `gorm:"not null;default:300" json:"coffee_break_max_sec"`

	QuietHoursStart string `gorm:"type:time;not null;default:'20:00'" json:"quiet_hours_start"`
	QuietHoursEnd   string `gorm:"type:time;not null;default:'08:00'" json:"quiet_hours_end"`
	RespectTimezone bool   `gorm:"not null;default:true" json:"respect_timezone"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ThrottleProfile) TableName() string { return "cmp_throttle_profiles" }
