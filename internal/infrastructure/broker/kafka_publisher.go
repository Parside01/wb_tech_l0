package broker

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
	"wb_tech_l0/internal/infrastructure/config"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

type KafkaMessage = kafka.Message

func NewKafkaPublisher() *KafkaPublisher {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.C.KafkaConfig.Brokers...),
		Topic:        config.C.KafkaConfig.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchTimeout: 100 * time.Millisecond,
	}
	return &KafkaPublisher{
		writer: writer,
	}
}

func (k *KafkaPublisher) PublishMessages(ctx context.Context, messages ...KafkaMessage) error {
	if err := k.writer.WriteMessages(ctx, messages...); err != nil {
		return err
	}
	return nil
}

func (k *KafkaPublisher) Close() error {
	k.writer.Stats()
	return k.writer.Close()
}
