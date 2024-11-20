package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"wb_tech_l0/internal/entity"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *entity.Order) error
	GetAll(ctx context.Context) (map[string]*entity.Order, error)
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	repo := &orderRepository{db: db}
	repo.Init()
	return repo
}

func (r *orderRepository) Init() {
	r.db.MustExec(scheme)
}

// TODO: Split to several functions.
func (r *orderRepository) SaveOrder(ctx context.Context, order *entity.Order) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	if err := r.saveDelivery(ctx, tx, order); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, savePayments,
		order.Payment.PaymentID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	); err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		if _, err := tx.ExecContext(ctx, saveItems,
			item.OrderID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NMID,
			item.Brand,
			item.Status,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	if _, err := tx.ExecContext(ctx, saveOrder,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Delivery.DeliveryID,
		order.Payment.PaymentID,
		order.Locate,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SMID,
		order.DateCreated,
		order.OofShard,
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) GetAll(ctx context.Context) (map[string]*entity.Order, error) {
	sql_orders := []*entity.SQLOrder{}
	if err := r.db.SelectContext(ctx, &sql_orders, getAllOrders); err != nil {
		return nil, err
	}

	orders := make(map[string]*entity.Order)
	for _, o := range sql_orders {
		delivery := &entity.Delivery{}
		payment := &entity.Payment{}
		items := []*entity.Item{}

		if err := r.db.GetContext(ctx, delivery, getDeliveryByID, o.DeliveryID); err != nil {
			return nil, fmt.Errorf("GetAll: GetDelivery: %s", err.Error())
		}
		if err := r.db.GetContext(ctx, payment, getPaymentsByID, o.PaymentID); err != nil {
			return nil, fmt.Errorf("GetAll: GetPayments: %s", err.Error())
		}
		if err := r.db.SelectContext(ctx, &items, getItemsByOrderID, o.OrderUID); err != nil {
			return nil, fmt.Errorf("GetAll: GetItems: %s", err.Error())
		}

		order := &entity.Order{}
		o.CopyToOrder(order)
		order.Delivery = delivery
		order.Payment = payment
		order.Items = items
		orders[o.OrderUID] = order
	}

	return orders, nil
}

func (r *orderRepository) saveDelivery(ctx context.Context, tx *sqlx.Tx, order *entity.Order) error {
	if _, err := tx.ExecContext(ctx, saveDelivery,
		order.Delivery.DeliveryID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
