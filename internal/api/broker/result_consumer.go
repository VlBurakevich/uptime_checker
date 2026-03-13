package broker

import (
	"context"
	"encoding/json"
	"log"
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
			log.Printf("kafka riding error: %v", err)
			continue
		}

		var result dto.SiteCheckResult
		if err := json.Unmarshal(m.Value, &result); err != nil {
			log.Printf("deserialization error: %v", err)
		}

		if err := handler(ctx, result); err != nil {
			log.Printf("handler error: %v", err)
		}
	}
}

func (c *ResultConsumer) Close() error {
	return c.reader.Close()
}
