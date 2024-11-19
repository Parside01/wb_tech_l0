package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"wb_tech_l0/internal/entity"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *entity.Order) error
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *entity.Order) error {
	return nil
}
