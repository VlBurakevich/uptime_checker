package broker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
	"uptime-checker/internal/shared/dto"

	"github.com/segmentio/kafka-go"
)

type TaskConsumer struct {
	reader     *kafka.Reader
	MaxTaskAge time.Duration
}

func NewTaskConsumer(addr, topic, groupID string, maxTaskAge time.Duration) *TaskConsumer {
	return &TaskConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{addr},
			Topic:   topic,
			GroupID: groupID,
		}),
		MaxTaskAge: maxTaskAge,
	}
}

func (c *TaskConsumer) Start(ctx context.Context, handler func(context.Context, dto.SiteCheckTask) error) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			slog.Error("kafka riding error", "error", err)
			continue
		}

		var task dto.SiteCheckTask
		if err := json.Unmarshal(m.Value, &task); err != nil {
			slog.Error("deserialization error", "error", err)
			continue
		}

		if !task.CreatedAt.IsZero() && time.Since(task.CreatedAt) > c.MaxTaskAge {
			slog.Warn("skipping state task",
				"site_id", task.SiteID,
				"age", time.Since(task.CreatedAt).Round(time.Second),
			)
			continue
		}

		if err := handler(ctx, task); err != nil {
			slog.Error("handler error", "error", err)
		}
	}
}

func (c *TaskConsumer) Close() error {
	return c.reader.Close()
}
