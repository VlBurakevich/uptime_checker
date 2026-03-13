package broker

import (
	"context"
	"encoding/json"
	"log/slog"
	"uptime-checker/internal/shared/dto"

	"github.com/segmentio/kafka-go"
)

type ResultConsumer struct {
	reader *kafka.Reader
}

func NewResultConsumer(addr, topic, groupID string) *ResultConsumer {
	return &ResultConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{addr},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *ResultConsumer) Start(ctx context.Context, handler func(context.Context, dto.SiteCheckResult) error) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			slog.Error("kafka riding error", "error", err)
			continue
		}

		var result dto.SiteCheckResult
		if err := json.Unmarshal(m.Value, &result); err != nil {
			slog.Error("deserialization error", "error", err)
			continue
		}

		if err := handler(ctx, result); err != nil {
			slog.Error("handler error", "error", err)
		}
	}
}

func (c *ResultConsumer) Close() error {
	return c.reader.Close()
}
