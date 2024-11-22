package broker

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	brokers []string
	reader  *kafka.Reader
}

// TODO: Если сейчас что-то упадет в процессе, то все поломается и потеряется.
// Надо подумать чего с этим сделать.
func NewKafkaConsumer(brokers []string, topic string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{
		brokers: brokers,
		reader:  reader,
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

	if err := c.reader.CommitMessages(ctx, message); err != nil {
		return fmt.Errorf("Error committing message to Kafka: %s", err)
	}
	return nil
}
