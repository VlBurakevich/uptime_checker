package broker

import (
	"context"
	"encoding/json"
	"uptime-checker/internal/dto"

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
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(task.SiteID.String()),
		Value: payload,
	})
}
