package app

import (
	"uptime-checker/internal/watcher/broker"
	"uptime-checker/internal/watcher/checker"
	"uptime-checker/internal/watcher/config"
)

func New(cfg *config.Config) (*WatcherApp, func(), error) {
	taskConsumer := broker.NewTaskConsumer(cfg.KafkaBroker, cfg.TopicSiteTask, "watcher-group")
	resultProducer := broker.NewResultProducer(cfg.KafkaBroker, cfg.TopicCheckResult)
	pinger := checker.NewPinger(cfg.HTTPTimeout)

	application := NewWatcherApp(taskConsumer, resultProducer, pinger)

	cleanup := func() {
		_ = taskConsumer.Close()
		_ = resultProducer.Close()
	}

	return application, cleanup, nil
}
