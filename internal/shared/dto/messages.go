package dto

import (
	"time"

	"github.com/google/uuid"
)

type SiteCheckTask struct {
	SiteID    uuid.UUID `json:"site_id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

type SiteCheckResult struct {
	SiteID     uuid.UUID `json:"site_id"`
	StatusCode int       `json:"status_code"`
	LatencyMs  int64     `json:"latency_ms"`
	IsUp       bool      `json:"is_up"`
	CheckedAt  time.Time `json:"checked_at"`
	Error      string    `json:"error,omitempty"`
}
