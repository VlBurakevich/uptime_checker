package checker

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"
	"uptime-checker/internal/shared/dto"
)

type Pinger struct {
	client *http.Client
}

func NewPinger(timeout time.Duration) *Pinger {
	return &Pinger{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (p *Pinger) Ping(ctx context.Context, task dto.SiteCheckTask) dto.SiteCheckResult {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", task.URL, nil)

	if err != nil {
		return dto.SiteCheckResult{
			SiteID:    task.SiteID,
			CheckedAt: time.Now(),
			IsUp:      false,
			Error:     "failed to create request: " + err.Error(),
		}
	}

	req.Header.Set("User-Agent", "UptimeWatcher/1.0 (University Project; Vlad's Bot)")

	resp, err := p.client.Do(req)

	latency := time.Since(start).Milliseconds()

	result := dto.SiteCheckResult{
		SiteID:    task.SiteID,
		CheckedAt: start,
		LatencyMs: latency,
	}

	if err != nil {
		result.IsUp = false
		result.Error = err.Error()
		result.StatusCode = -1
		return result
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}(resp.Body)

	_, _ = io.Copy(io.Discard, resp.Body)

	result.StatusCode = resp.StatusCode
	result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 300

	return result
}
