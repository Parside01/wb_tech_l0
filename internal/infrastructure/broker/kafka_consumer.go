package broker

import (
	"context"
	"github.com/segmentio/kafka-go"
	"wb_tech_l0/internal/infrastructure/config"
)

type KafkaConsumer interface {
	ConsumeMessage(ctx context.Context) (KafkaMessage, error)
	Close() error
}

type kafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer() KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  config.C.KafkaConfig.Brokers,
		Topic:    config.C.KafkaConfig.Topic,
		GroupID:  config.C.KafkaConfig.GroupID,
		MinBytes: 1,
		MaxBytes: config.C.KafkaConfig.MaxBytes,
	})
	return &kafkaConsumer{
		reader: reader,
	}
}

func (c *kafkaConsumer) ConsumeMessage(ctx context.Context) (KafkaMessage, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *kafkaConsumer) Close() error {
	return c.reader.Close()
}
