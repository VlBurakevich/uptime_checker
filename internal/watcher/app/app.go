package app

import (
	"context"
	"log/slog"
	"uptime-checker/internal/shared/dto"
	"uptime-checker/internal/watcher/broker"
	"uptime-checker/internal/watcher/checker"
	"uptime-checker/internal/watcher/dispatcher"
)

type WatcherApp struct {
	consumer   *broker.TaskConsumer
	producer   *broker.ResultProducer
	pinger     *checker.Pinger
	dispatcher *dispatcher.Dispatcher
	monitor    *dispatcher.AdaptiveMonitor
}

func NewWatcherApp(
	cons *broker.TaskConsumer,
	prod *broker.ResultProducer,
	ping *checker.Pinger,
	disp *dispatcher.Dispatcher,
	monitor *dispatcher.AdaptiveMonitor,
) *WatcherApp {
	return &WatcherApp{
		consumer:   cons,
		producer:   prod,
		pinger:     ping,
		dispatcher: disp,
		monitor:    monitor,
	}
}

func (a *WatcherApp) Run(ctx context.Context) error {
	go a.monitor.Run(ctx)

	return a.consumer.Start(ctx, a.handleTask)
}

func (a *WatcherApp) handleTask(ctx context.Context, task dto.SiteCheckTask) error {
	a.dispatcher.Execute(ctx, func() {
		result := a.pinger.Ping(ctx, task)
		if err := a.producer.PublishResult(ctx, result); err != nil {
			slog.Error("failed to publish result",
				"site_id", task.SiteID,
				"url", task.URL,
				"error", err,
			)
		}
	})

	return nil
}
