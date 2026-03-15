package app

import (
	"uptime-checker/internal/watcher/broker"
	"uptime-checker/internal/watcher/checker"
	"uptime-checker/internal/watcher/config"
	"uptime-checker/internal/watcher/dispatcher"
)

func New(cfg *config.Config) (*WatcherApp, func(), error) {
	taskConsumer := broker.NewTaskConsumer(
		cfg.Kafka.Broker,
		cfg.Kafka.TopicTasks,
		cfg.Kafka.GroupId,
		cfg.Watcher.MaxTaskAge,
	)

	resultProducer := broker.NewResultProducer(cfg.Kafka.Broker, cfg.Kafka.TopicResults)
	pinger := checker.NewPinger(cfg.Watcher.HTTPTimeout)
	disp := dispatcher.New(cfg.Watcher.MinConcurrency, cfg.Watcher.MaxConcurrency)

	monitor := dispatcher.NewAdaptiveMonitor(
		disp,
		cfg.Watcher.Monitor.CpuLow,
		cfg.Watcher.Monitor.CpuHigh,
		cfg.Watcher.Monitor.Increment,
		cfg.Watcher.Monitor.Throttle,
		cfg.Watcher.Monitor.Interval,
	)

	application := NewWatcherApp(taskConsumer, resultProducer, pinger, disp, monitor)

	cleanup := func() {
		_ = taskConsumer.Close()
		disp.Wait()
		_ = resultProducer.Close()
	}

	return application, cleanup, nil
}
