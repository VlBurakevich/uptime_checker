package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"uptime-checker/internal/api/app"
	"uptime-checker/internal/api/config"
)

func main() {
	cfg := config.Load()

	application, cleanup, err := app.New(cfg)
	if err != nil {
		slog.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("App is running")
	if err := application.Run(ctx); err != nil {
		slog.Error("App crashed", "error", err)
		os.Exit(1)
	}
}
