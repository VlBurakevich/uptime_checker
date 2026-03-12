package broker

import (
	"context"
	"encoding/json"
	"log"
	"uptime-checker/internal/dto"

	"github.com/segmentio/kafka-go"
)

type TaskConsumer struct {
	reader *kafka.Reader
}

func NewTaskConsumer(addr, topic, groupID string) *TaskConsumer {
	return &TaskConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{addr},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *TaskConsumer) Start(ctx context.Context, handler func(context.Context, dto.SiteCheckTask)) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("kafka riding error: %v", err)
			continue
		}

		var task dto.SiteCheckTask
		if err := json.Unmarshal(m.Value, &task); err != nil {
			log.Printf("deserialization error: %v", err)
			continue
		}

		handler(ctx, task)
	}
}

func (c *TaskConsumer) Close() error {
	return c.reader.Close()
}
