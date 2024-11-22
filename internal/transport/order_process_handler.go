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
	logger   *zap.Logger
	consumer *broker.KafkaConsumer
}

func NewOrderProcessHandler(service service.OrderService, consumer *broker.KafkaConsumer, logger *zap.Logger) *OrderProcessHandler {
	return &OrderProcessHandler{
		service:  service,
		logger:   logger,
		consumer: consumer,
	}
}

func (h *OrderProcessHandler) Start(ctx context.Context) error {
	for {
		if err := h.consumer.ConsumeMessage(ctx, h.loggingMiddleware(h.handleMessage)); err != nil {
			return err
		}
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
		h.logger.Info("Received message from kafka", zap.String("key", string(message.Key)), zap.String("value", string(message.Value)))

		err := next(ctx, message)

		duration := time.Since(start)
		if err != nil {
			h.logger.Error("Error processing message", zap.String("key", string(message.Key)), zap.Error(err), zap.Duration("processing_time", duration))
		} else {
			h.logger.Info("Message processed successfully", zap.String("key", string(message.Key)), zap.Duration("processing_time", duration))
		}

		return err
	}
}
