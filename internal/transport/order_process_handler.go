package transport

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
	"wb_tech_l0/internal/entity"
	"wb_tech_l0/internal/infrastructure/broker"
	"wb_tech_l0/internal/service"
)

type OrderProcessHandler struct {
	service  service.OrderService
	consumer *broker.KafkaConsumer
}

func NewOrderProcessHandler(service service.OrderService, consumer *broker.KafkaConsumer) *OrderProcessHandler {
	return &OrderProcessHandler{
		service:  service,
		consumer: consumer,
	}
}

func (h *OrderProcessHandler) Start(ctx context.Context) error {
	for {
		// Все ошибки и так будут логироваться, так что нет смысла что-то с ней делать.
		_ = h.consumer.ConsumeMessage(ctx, h.loggingMiddleware(h.handleMessage))
	}
}

func (h *OrderProcessHandler) handleMessage(ctx context.Context, message kafka.Message) error {
	var order *entity.Order
	if err := json.Unmarshal(message.Value, &order); err != nil {
		return err
	}
	if err := h.service.SaveOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (h *OrderProcessHandler) loggingMiddleware(next broker.KafkaMessageHandler) broker.KafkaMessageHandler {
	return func(ctx context.Context, message kafka.Message) error {
		start := time.Now()

		zap.L().Info("Received message from kafka", zap.String("key", string(message.Key)))

		err := next(ctx, message)

		duration := time.Since(start)
		if err != nil {
			zap.L().Error("Error processing message", zap.String("key", string(message.Key)), zap.Error(err), zap.Duration("processing_time", duration))
		} else {
			zap.L().Info("Message processed successfully", zap.String("key", string(message.Key)), zap.Duration("processing_time", duration))
		}

		return err
	}
}
