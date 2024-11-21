package broker

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

type KafkaConsumer struct {
	brokers []string
	reader  *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})

	return &KafkaConsumer{
		brokers: brokers,
		reader:  reader,
	}
}

func (c *KafkaConsumer) Consume(ctx context.Context, handler func(ctx context.Context, message kafka.Message) error) error {
	message, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("Error reading message from Kafka: %s", err)
	}
	if err := handler(ctx, message); err != nil {
		return fmt.Errorf("Error processing message from Kafka: %s", err)
	}
	return nil
}
