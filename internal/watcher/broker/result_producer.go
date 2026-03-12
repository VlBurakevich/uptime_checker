package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"uptime-checker/internal/shared/dto"

	"github.com/segmentio/kafka-go"
)

type ResultProducer struct {
	writer *kafka.Writer
}

func NewResultProducer(addr, topic string) *ResultProducer {
	return &ResultProducer{
		writer: &kafka.Writer{
			Addr:        kafka.TCP(addr),
			Topic:       topic,
			Balancer:    &kafka.LeastBytes{},
			Compression: kafka.Snappy,
		},
	}
}

func (p *ResultProducer) PublishResult(ctx context.Context, result dto.SiteCheckResult) error {
	payload, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(result.SiteID.String()),
		Value: payload,
	})

	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (p *ResultProducer) Close() error {
	return p.writer.Close()
}
