package broker

import (
	"context"
	"github.com/segmentio/kafka-go"
	"wb_tech_l0/internal/infrastructure/config"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

type KafkaMessageHandler = func(ctx context.Context, message KafkaMessage) error

// TODO: Если сейчас что-то упадет в процессе, то все поломается и потеряется.
// Надо подумать чего с этим сделать.
func NewKafkaConsumer() *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  config.C.KafkaConfig.Brokers,
		Topic:    config.C.KafkaConfig.Topic,
		GroupID:  config.C.KafkaConfig.GroupID,
		MinBytes: 1,
		MaxBytes: config.C.KafkaConfig.MaxBytes,
	})
	return &KafkaConsumer{
		reader: reader,
	}
}

func (c *KafkaConsumer) ConsumeMessage(ctx context.Context) (KafkaMessage, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
