package transport

import (
	"context"
	"encoding/json"
	"github.com/gammazero/workerpool"
	"go.uber.org/zap"
	"sync"
	"time"
	"wb_tech_l0/internal/entity"
	"wb_tech_l0/internal/infrastructure/broker"
	"wb_tech_l0/internal/service"
)

type OrderProcessHandler struct {
	service  service.OrderService
	consumer broker.KafkaConsumer
	workers  *workerpool.WorkerPool
}

type metric struct {
	ProcessedReqCount int64
	mutex             sync.Mutex
}

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
		start := time.Now()

		if err := h.processMessage(ctx, message); err != nil {
			zap.L().Error("Error processing message", zap.String("key", string(message.Key)), zap.Error(err), zap.Duration("processing_time", time.Since(start)))
		} else {
			zap.L().Info("KafkaMessage processed successfully", zap.String("key", string(message.Key)), zap.Duration("processing_time", time.Since(start)))

			if err := h.consumer.CommitMessages(ctx, message); err != nil {
				zap.L().Error("Failed commit message", zap.String("key", string(message.Key)), zap.Error(err))
			}
		}
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
