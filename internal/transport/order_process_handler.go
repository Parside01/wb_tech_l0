package transport

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"sync"
	"time"
	"wb_tech_l0/internal/entity"
	"wb_tech_l0/internal/infrastructure/broker"
	"wb_tech_l0/internal/service"
)

type OrderProcessHandler struct {
	service      service.OrderService
	consumer     *broker.KafkaConsumer
	workersCount int
	metrics      metric
}

type metric struct {
	ProcessedReqCount int64
	mutex             sync.Mutex
}

func NewOrderProcessHandler(service service.OrderService, consumer *broker.KafkaConsumer) *OrderProcessHandler {
	return &OrderProcessHandler{
		service:      service,
		consumer:     consumer,
		workersCount: 100,
		metrics: metric{
			ProcessedReqCount: 0,
			mutex:             sync.Mutex{},
		},
	}
}

func (h *OrderProcessHandler) Start(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	for i := 0; i < h.workersCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			h.startWorker(ctx)
		}()
	}
	go h.startMetricsLogger()
	wg.Wait()
	return nil
}

func (h *OrderProcessHandler) startMetricsLogger() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	average := int64(0)
	for range ticker.C {
		h.metrics.mutex.Lock()
		average = max(average, h.metrics.ProcessedReqCount)
		h.metrics.ProcessedReqCount = 0
		h.metrics.mutex.Unlock()

		zap.L().Warn("", zap.Int64("RPS", average))
	}
}

func (h *OrderProcessHandler) handleMessage(ctx context.Context, message broker.KafkaMessage) error {
	var order *entity.Order
	if err := json.Unmarshal(message.Value, &order); err != nil {
		return err
	}
	if err := h.service.SaveOrder(ctx, order); err != nil {
		return err
	}

	h.metrics.mutex.Lock()
	h.metrics.ProcessedReqCount++
	h.metrics.mutex.Unlock()

	return nil
}

func (h *OrderProcessHandler) startWorker(ctx context.Context) {
	for {
		message, err := h.consumer.ConsumeMessage(ctx)
		if err != nil {
			zap.L().Error("Failed to consume message", zap.Error(err))
			continue
		}

		start := time.Now()

		zap.L().Info("Received message from kafka", zap.String("key", string(message.Key)))

		if err := h.handleMessage(ctx, message); err != nil {
			zap.L().Error("Error processing message", zap.String("key", string(message.Key)), zap.Error(err), zap.Duration("processing_time", time.Since(start)))
		} else {
			zap.L().Info("KafkaMessage processed successfully", zap.String("key", string(message.Key)), zap.Duration("processing_time", time.Since(start)))
		}

	}
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
			zap.L().Info("KafkaMessage processed successfully", zap.String("key", string(message.Key)), zap.Duration("processing_time", duration))
		}

		return err
	}
}
