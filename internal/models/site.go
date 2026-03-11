package models

import (
	"time"

	"github.com/google/uuid"
)

type Site struct {
	Base
	URL      string        `gorm:"not null" json:"url"`
	Name     string        `json:"name"`
	Interval time.Duration `gorm:"default:60s" json:"interval"`
	IsActive bool          `gorm:"default:true" json:"is_active"`
	UserId   uuid.UUID     `gorm:"type:uuid" json:"user_id"`
	Checks   []SiteCheck   `gorm:"foreignKey:SiteID" json:"checks"`
}

type SiteCheck struct {
	Identity
	SiteID     uuid.UUID     `gorm:"type:uuid;index"`
	StatusCode int           `json:"status_code"`
	Latency    time.Duration `json:"latency"`
	IsUp       bool          `json:"is_up"`
	CheckedAt  time.Time     `gorm:"autoCreateTime" json:"checked_at"`
}
