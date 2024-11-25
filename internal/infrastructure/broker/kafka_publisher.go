package broker

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
	"wb_tech_l0/internal/infrastructure/config"
)

type KafkaMessage = kafka.Message

type KafkaPublisher interface {
	PublishMessages(ctx context.Context, messages ...KafkaMessage) error
	Close() error
}

type kafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher() KafkaPublisher {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.C.KafkaConfig.Brokers...),
		Topic:        config.C.KafkaConfig.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
	}
	return &kafkaPublisher{
		writer: writer,
	}
}

func (k *kafkaPublisher) PublishMessages(ctx context.Context, messages ...KafkaMessage) error {
	if err := k.writer.WriteMessages(ctx, messages...); err != nil {
		return err
	}
	return nil
}

func (k *kafkaPublisher) Close() error {
	return k.writer.Close()
}
