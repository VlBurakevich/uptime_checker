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

	result := dto.SiteCheckResult{
		SiteID:    task.SiteID,
		CheckedAt: start,
	}
	if err != nil {
		result.IsUp = false
		result.Error = err.Error()
		return result
	}

	resp, err := p.client.Do(req)
	result.LatencyMs = time.Since(start).Milliseconds()

	if err != nil {
		result.IsUp = false
		result.Error = err.Error()
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
