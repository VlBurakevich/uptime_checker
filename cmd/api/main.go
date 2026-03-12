package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uptime-checker/internal/api"
	"uptime-checker/internal/api/broker"
	"uptime-checker/internal/api/config"
	"uptime-checker/internal/api/database"
	"uptime-checker/internal/api/services"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatal(err)
	}

	taskProducer := broker.NewTaskProducer(cfg.Kafka.Addr, cfg.Kafka.TopicTasks)

	scheduler := services.NewScheduler(db, taskProducer, cfg.Scheduler.Limit)

	r := api.SetupRouter(db, config.GetJWTSecret())

	go scheduler.Run(ctx, cfg.Scheduler.Interval)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.AppPort)

	<-ctx.Done()
	log.Println("Gracefully shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped.")
}
