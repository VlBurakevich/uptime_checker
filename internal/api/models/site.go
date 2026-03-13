package models

import (
	"time"

	"github.com/google/uuid"
)

type Site struct {
	Base
	URL           string      `gorm:"not null" json:"url"`
	Name          string      `json:"name"`
	Interval      int         `gorm:"default:60" json:"interval"`
	LastCheckedAt time.Time   `gorm:"type:timestamptz" json:"last_checked_at"`
	NextCheckedAt time.Time   `gorm:"index:idx_active_next;type:timestamptz" json:"next_checked_at"`
	IsActive      bool        `gorm:"index:idx_active_next;default:true" json:"is_active"`
	UserID        uuid.UUID   `gorm:"type:uuid" json:"user_id"`
	Checks        []SiteCheck `gorm:"foreignKey:SiteID" json:"checks"`
}

type SiteCheck struct {
	Identity
	SiteID     uuid.UUID `gorm:"type:uuid;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StatusCode int       `json:"status_code"`
	LatencyMs  int64     `json:"latency"`
	IsUp       bool      `json:"is_up"`
	CheckedAt  time.Time `gorm:"autoCreateTime" json:"checked_at"`
}
