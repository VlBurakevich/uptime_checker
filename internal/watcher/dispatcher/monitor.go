package dispatcher

import (
	"context"
	"log/slog"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

func (d *Dispatcher) RunAdaptiveMonitor(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	slog.Info("Adaptive monitor started",
		"min", d.minLimit,
		"max", d.maxLimit)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Adaptive monitor stopped")
			return
		case <-ticker.C:
			percentages, err := cpu.PercentWithContext(ctx, time.Second, false)
			if err != nil {
				slog.Error("Failed to get CPU usage", "error", err)
			}

			if len(percentages) == 0 {
				continue
			}

			usage := percentages[0]
			currentLimit := d.targetConcurrency
			var nextLimit int32

			switch {
			case usage < 50:
				nextLimit = currentLimit + 5
				slog.Debug("increasing limit", "cpu", usage, "new_limit", nextLimit)
			case usage > 80:
				nextLimit = int32(float64(currentLimit) * 0.8)
				slog.Warn("Cpu overload detected, throttling", "cpu", usage, "new_limit", nextLimit)
			default:
				continue
			}
			d.SetLimit(nextLimit)
		}
	}
}
