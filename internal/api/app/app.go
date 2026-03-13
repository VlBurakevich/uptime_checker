package app

import (
	"context"
	"errors"
	"log"
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
			log.Printf("Result consumer stopped with error: %v", err)
		}
	}()

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()
	return a.Stop()
}

func (a *App) Stop() error {
	return a.server.Shutdown(context.Background())
}
