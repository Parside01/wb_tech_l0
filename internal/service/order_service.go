package service

import (
	"context"
	"wb_tech_l0/internal/entity"
	"wb_tech_l0/internal/infrastructure/cache"
	"wb_tech_l0/internal/repository"
)

const (
	IDsSize = 32
)

type OrderService interface {
	SaveOrder(ctx context.Context, order *entity.Order) error
	GetOrderById(ctx context.Context, id string) (*entity.Order, error)
}

type orderService struct {
	repo  repository.OrderRepository
	cache cache.Cache
}

// TODO: Надо бы наверное тесты написать..
func NewOrderService(repo repository.OrderRepository, cache cache.Cache) OrderService {
	return &orderService{
		repo:  repo,
		cache: cache,
	}
}

// TODO: Не знаю на самом деле, кешировать это точно обязанность сервиса?
// TODO: Нам нужно учитывать случай, когда пришел дубликат существуещего заказа?
func (s *orderService) SaveOrder(ctx context.Context, order *entity.Order) error {
	order.Delivery.OrderID = order.OrderUID
	order.Payment.OrderID = order.OrderUID
	for _, i := range order.Items {
		i.OrderID = order.OrderUID
	}
	if err := s.repo.Save(ctx, order); err != nil {
		return err
	}
	s.cache.Set(order.Key(), order)
	return nil
}

func (s *orderService) GetOrderById(ctx context.Context, id string) (*entity.Order, error) {
	if cachedOrder, ok := s.cache.Get(id); ok {
		return cachedOrder.(*entity.Order), nil
	}
	return &entity.Order{}, ErrNoOrderInCache
}
