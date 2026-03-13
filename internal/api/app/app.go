package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"uptime-checker/internal/api/broker"
	"uptime-checker/internal/api/config"
	"uptime-checker/internal/api/services"

	"gorm.io/gorm"
)

type App struct {
	cfg            *config.Config
	db             *gorm.DB
	scheduler      *services.Scheduler
	resultConsumer *broker.ResultConsumer
	siteService    *services.SiteService
	server         *http.Server
}

func (a *App) Run(ctx context.Context) error {
	go a.scheduler.Run(ctx, a.cfg.Scheduler.Interval)

	go func() {
		if err := a.resultConsumer.Start(ctx, a.siteService.HandleCheckResult); err != nil {
			slog.Error("Result consumer stopped", "error", err)
		}
	}()

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failed", "error", err)
		}
	}()

	<-ctx.Done()
	return a.Stop()
}

func (a *App) Stop() error {
	return a.server.Shutdown(context.Background())
}
