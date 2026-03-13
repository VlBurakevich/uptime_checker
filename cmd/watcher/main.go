package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"uptime-checker/internal/watcher/app"
	"uptime-checker/internal/watcher/config"
)

func main() {
	cfg := config.Load()

	application, cleanup, err := app.New(cfg)
	if err != nil {
		slog.Error("Failed to initialize app", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("Watcher is running")
	if err := application.Run(ctx); err != nil {
		slog.Error("App crashed", "error", err)
		os.Exit(1)
	}
}
