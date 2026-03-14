package app

import (
	"context"
	"log"
	"uptime-checker/internal/shared/dto"
	"uptime-checker/internal/watcher/broker"
	"uptime-checker/internal/watcher/checker"
)

type WatcherApp struct {
	consumer *broker.TaskConsumer
	producer *broker.ResultProducer
	pinger   *checker.Pinger
}

func NewWatcherApp(cons *broker.TaskConsumer, prod *broker.ResultProducer, ping *checker.Pinger) *WatcherApp {
	return &WatcherApp{
		consumer: cons,
		producer: prod,
		pinger:   ping,
	}
}

func (a *WatcherApp) Run(ctx context.Context) error {
	return a.consumer.Start(ctx, a.handleTask)
}

func (a *WatcherApp) handleTask(ctx context.Context, task dto.SiteCheckTask) error {
	go func() {
		result := a.pinger.Ping(ctx, task)

		if err := a.producer.PublishResult(ctx, result); err != nil {
			log.Printf("Error publishing result for site %s: %v", task.SiteID, err)
		}
	}()
	return nil
}
