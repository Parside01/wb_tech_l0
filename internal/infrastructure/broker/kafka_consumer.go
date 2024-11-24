package broker

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"wb_tech_l0/internal/infrastructure/config"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

type KafkaMessageHandler = func(ctx context.Context, msg kafka.Message) error

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

// TODO: Можно добавить ретраи, но потом если так и не получилось обработать пихать это куда-нибудь.
// TODO: Щас просто надо mvp сдеалать, а потом если что worker-pool добавить.
func (c *KafkaConsumer) ConsumeMessage(ctx context.Context, handler func(ctx context.Context, message kafka.Message) error) error {
	message, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("Error reading message from Kafka: %s", err)
	}
	if err := handler(ctx, message); err != nil {
		return fmt.Errorf("Error processing message from Kafka: %s", err)
	}
	return nil
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
