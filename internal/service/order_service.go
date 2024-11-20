package service

import (
	"context"
	"github.com/twharmon/gouid"
	"wb_tech_l0/internal/entity"
	"wb_tech_l0/internal/repository"
)

const (
	IDsSize = 32
)

type OrderService interface {
	SaveOrder(ctx context.Context, order *entity.Order) error
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{
		repo: repo,
	}
}

func (s *orderService) SaveOrder(ctx context.Context, order *entity.Order) error {
	deliveryId := gouid.Bytes(IDsSize).String()
	paymentId := gouid.Bytes(IDsSize).String()

	order.Delivery.DeliveryID = deliveryId
	order.Payment.PaymentID = paymentId
	for _, i := range order.Items {
		i.OrderID = order.OrderUID
	}

	return s.repo.Save(ctx, order)
}
