package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"uptime-checker/internal/shared/dto"

	"github.com/segmentio/kafka-go"
)

type TaskProducer struct {
	writer *kafka.Writer
}

func NewTaskProducer(addr, topic string) *TaskProducer {
	return &TaskProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(addr),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *TaskProducer) PublishTask(ctx context.Context, task dto.SiteCheckTask) error {
	payload, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(task.SiteID.String()),
		Value: payload,
	})

	if err != nil {
		return fmt.Errorf("failed to publish task: %w", err)
	}

	return nil
}

func (p *TaskProducer) Close() error {
	return p.writer.Close()
}
