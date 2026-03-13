package app

import (
	"log"
	"net/http"
	"uptime-checker/internal/api"
	"uptime-checker/internal/api/broker"
	"uptime-checker/internal/api/config"
	"uptime-checker/internal/api/database"
	"uptime-checker/internal/api/services"
)

func New(cfg *config.Config) (*App, func(), error) {
	db, err := database.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatal("DB init error: ", err)
		return nil, nil, err
	}

	taskProducer := broker.NewTaskProducer(cfg.Kafka.KafkaBroker, cfg.Kafka.TopicTasks)

	siteService := &services.SiteService{DB: db}
	scheduler := services.NewScheduler(db, taskProducer, cfg.Scheduler.Limit)

	resultConsumer := broker.NewResultConsumer(cfg.Kafka.KafkaBroker, cfg.Kafka.TopicResults, "api-group")

	r := api.SetupRouter(db, cfg.JWTSecret)

	application := &App{
		cfg:            cfg,
		db:             db,
		scheduler:      scheduler,
		resultConsumer: resultConsumer,
		siteService:    siteService,
		server: &http.Server{
			Addr:    ":" + cfg.AppPort,
			Handler: r,
		},
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		_ = taskProducer.Close()
		_ = resultConsumer.Close()
	}

	return application, cleanup, nil
}
