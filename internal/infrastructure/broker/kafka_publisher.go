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

func NewKafkaPublisher() *KafkaPublisher {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.C.KafkaConfig.Brokers...),
		Topic:        config.C.KafkaConfig.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchTimeout: 5 * time.Millisecond,
	}
	return &KafkaPublisher{
		writer: writer,
	}
}

func (k *KafkaPublisher) PublishMessage(ctx context.Context, key, data []byte) error {
	message := kafka.Message{
		Key:   key,
		Value: data,
	}
	if err := k.writer.WriteMessages(ctx, message); err != nil {
		return err
	}
	return nil
}

func (k *KafkaPublisher) Close() error {
	k.writer.Stats()
	return k.writer.Close()
}