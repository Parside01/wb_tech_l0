package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"wb_tech_l0/internal/entity"
)

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
	GetAll(ctx context.Context) (map[string]*entity.Order, error)
}

type orderRepository struct {
	db       *sqlx.DB
	payments PaymentRepository
	delivery DeliveryRepository
	items    ItemRepository
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	repo := &orderRepository{
		db:       db,
		items:    NewItemRepository(db),
		delivery: NewDeliveryRepository(db),
		payments: NewPaymentRepository(db),
	}
	repo.Init()
	return repo
}

func (r *orderRepository) Init() {
	r.db.MustExec(scheme)
}

func (r *orderRepository) Save(ctx context.Context, order *entity.Order) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	if err := r.delivery.SaveTX(ctx, tx, order.Delivery); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := r.payments.SaveTX(ctx, tx, order.Payment); err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		if err := r.items.SaveTX(ctx, tx, item); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := r.saveOrderTX(ctx, tx, order); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) GetAll(ctx context.Context) (map[string]*entity.Order, error) {
	orders := []*entity.Order{}
	if err := r.db.SelectContext(ctx, &orders, getAllOrders); err != nil {
		return nil, err
	}

	orders_map := make(map[string]*entity.Order)
	for _, order := range orders {
		delivery, err := r.delivery.GetByOrderID(ctx, order.OrderUID)
		if err != nil {
			return nil, err
		}

		payment, err := r.payments.GetByOrderID(ctx, order.OrderUID)
		if err != nil {
			return nil, err
		}

		items, err := r.items.GetAllByOrderID(ctx, order.OrderUID)
		if err != nil {
			return nil, err
		}

		order.Delivery = delivery
		order.Payment = payment
		order.Items = items
		orders_map[order.OrderUID] = order
	}

	return orders_map, nil
}

func (r *orderRepository) saveOrderTX(ctx context.Context, tx *sqlx.Tx, order *entity.Order) error {
	if _, err := tx.ExecContext(ctx, saveOrder,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locate,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SMID,
		order.DateCreated,
		order.OofShard,
	); err != nil {
		return err
	}
	return nil
}
