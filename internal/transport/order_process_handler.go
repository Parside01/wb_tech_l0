package transport

import (
	"context"
	"encoding/json"
	"github.com/gammazero/workerpool"
	"go.uber.org/zap"
	"time"
	"wb_tech_l0/internal/entity"
	"wb_tech_l0/internal/infrastructure/broker"
	"wb_tech_l0/internal/service"
	"wb_tech_l0/internal/transport/metrics"
)

type OrderProcessHandler struct {
	service  service.OrderService
	consumer broker.KafkaConsumer
	workers  *workerpool.WorkerPool
}

type Handler func(ctx context.Context, args ...interface{}) error
type Middleware func(Handler) Handler

func NewOrderProcessHandler(service service.OrderService, consumer broker.KafkaConsumer) *OrderProcessHandler {
	return &OrderProcessHandler{
		service:  service,
		consumer: consumer,
		workers:  workerpool.New(50),
	}
}

func (h *OrderProcessHandler) Start(ctx context.Context) error {
	for i := 0; i < h.workers.Size(); i++ {
		h.workers.Submit(func() {
			h.listenAndProcessMessages(ctx)
		})
	}
	return nil
}

func (h *OrderProcessHandler) Shutdown() error {
	h.workers.StopWait()
	err := h.consumer.Close()
	return err
}

func (h *OrderProcessHandler) listenAndProcessMessages(ctx context.Context) {
	for {
		message, err := h.consumer.FetchMessage(ctx)
		if err != nil {
			zap.L().Error("Failed to received message", zap.Error(err))
			return
		}

		zap.L().Info("Received message from kafka", zap.String("key", string(message.Key)))
		metrics.KafkaMessagesReceivedCount.WithLabelValues(message.Topic).Inc()

		start := time.Now()
		err = h.processMessage(ctx, message)
		duration := time.Since(start)

		if err != nil {
			zap.L().Error("Error processing message", zap.String("key", string(message.Key)), zap.Error(err), zap.Duration("processing_time", duration))
		} else {
			zap.L().Info("KafkaMessage processed successfully", zap.String("key", string(message.Key)), zap.Duration("processing_time", duration))

			if err = h.consumer.CommitMessages(ctx, message); err != nil {
				zap.L().Error("Failed commit message", zap.String("key", string(message.Key)), zap.Error(err))
			}
		}
		metrics.KafkaMessageProcessingDuration.WithLabelValues(message.Topic).Observe(duration.Seconds())
	}
}

func (h *OrderProcessHandler) processMessage(ctx context.Context, message broker.KafkaMessage) error {
	var order *entity.Order
	if err := json.Unmarshal(message.Value, &order); err != nil {
		return err
	}
	if err := h.service.SaveOrder(ctx, order); err != nil {
		return err
	}

	return nil
}
