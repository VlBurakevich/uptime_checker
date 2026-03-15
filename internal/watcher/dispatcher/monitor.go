package dispatcher

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type AdaptiveMonitor struct {
	dispatcher *Dispatcher

	lowThreshold  float64
	highThreshold float64
	increment     int32
	throttleRatio float64
	interval      time.Duration
}

func NewAdaptiveMonitor(
	d *Dispatcher,
	low, high float64,
	inc int32,
	ratio float64,
	interval time.Duration,
) *AdaptiveMonitor {
	return &AdaptiveMonitor{
		dispatcher:    d,
		lowThreshold:  low,
		highThreshold: high,
		increment:     inc,
		throttleRatio: ratio,
		interval:      interval,
	}
}

func (am *AdaptiveMonitor) Run(ctx context.Context) {
	ticker := time.NewTicker(am.interval)
	defer ticker.Stop()

	slog.Info("Adaptive monitor started",
		"min", am.dispatcher.minLimit,
		"max", am.dispatcher.maxLimit)

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

			usage := percentages[0]
			currentLimit := atomic.LoadInt32(&am.dispatcher.targetConcurrency)
			var nextLimit int32

			switch {
			case usage < am.lowThreshold:
				nextLimit = currentLimit + am.increment
				slog.Debug("increasing limit", "cpu", usage, "new_limit", nextLimit)
			case usage > am.highThreshold:
				nextLimit = int32(float64(currentLimit) * am.throttleRatio)
				slog.Warn("Cpu overload detected, throttling", "cpu", usage, "new_limit", nextLimit)
			default:
				continue
			}
			am.dispatcher.SetLimit(nextLimit)
		}
	}
}
