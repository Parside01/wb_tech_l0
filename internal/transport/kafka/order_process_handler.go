package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
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
		if err := h.consumer.ConsumeMessage(ctx, h.handleMessage); err != nil {
			return err
		}
	}
}

func (h *OrderProcessHandler) handleMessage(ctx context.Context, message kafka.Message) error {
	h.logger.Info("Received message from kafka", zap.String("key", string(message.Key)), zap.String("value", string(message.Value)))

	var order *entity.Order
	if err := json.Unmarshal(message.Value, &order); err != nil {
		h.logger.Error("Failed to unmarshal", zap.String("key", string(message.Key)), zap.String("value", string(message.Value)))
		return err
	}
	if err := h.service.SaveOrder(ctx, order); err != nil {
		h.logger.Error("Failed to save order", zap.Error(err))
	}

	h.logger.Info("Order saved successful", zap.String("key", string(message.Key)))
	return nil
}
