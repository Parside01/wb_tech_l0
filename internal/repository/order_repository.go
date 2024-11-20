package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"wb_tech_l0/internal/entity"
)

const (
	scheme = `
		CREATE TABLE IF NOT EXISTS delivery
		(
			delivery_id TEXT NOT NULL,
			name TEXT NOT NULL,
			phone TEXT NOT NULL,
			zip TEXT NOT NULL,
			city TEXT NOT NULL,
			address TEXT NOT NULL,
			region TEXT NOT NULL,
			email TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS payments
		(
		    payment_id TEXT NOT NULL,
			transaction TEXT NOT NULL,
			request_id TEXT NOT NULL,
			currency TEXT NOT NULL,
			provider TEXT NOT NULL,
			amount BIGINT NOT NULL,
			payment_dt BIGINT NOT NULL,
			bank TEXT NOT NULL,
			delivery_cost BIGINT NOT NULL,
			goods_total BIGINT NOT NULL,
			custom_fee BIGINT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS orders
		(
			order_uid TEXT NOT NULL,
			track_number TEXT NOT NULL,
			entry TEXT NOT NULL,
			delivery_id TEXT NOT NULL,
			payment_id TEXT NOT NULL,
			locate TEXT NOT NULL,
			internal_signature TEXT,
			customer_id TEXT NOT NULL,
			delivery_service TEXT NOT NULL,
			shard_key TEXT NOT NULL,
			sm_id TEXT NOT NULL,
			date_created TIMESTAMP,
			oof_shard TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS items
		(
			order_id TEXT NOT NULL,
			chrt_id BIGINT NOT NULL,
			track_number TEXT NOT NULL,
			price BIGINT NOT NULL,
			rid TEXT NOT NULL,
			name TEXT NOT NULL,
			sale BIGINT NOT NULL,
			size TEXT NOT NULL,
			total_price BIGINT NOT NULL,
			nm_id BIGINT NOT NULL,
			brand TEXT NOT NULL,
			status BIGINT NOT NULL
		);
	`

	saveOrder = `
		INSERT INTO orders
			(order_uid, track_number, entry, delivery_id, payment_id, locate, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);
	`

	saveDelivery = `
		INSERT INTO delivery
			(delivery_id, name, phone, zip, city, address, region, email)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8);
	`

	savePayments = `
		INSERT INTO payments
			(payment_id, transaction, request_id, currency, provider, amount,
			 payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`

	saveItems = `
		INSERT INTO items
			(order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
	`

	getAllOrders = `
		SELECT * FROM orders;
	`

	getDeliveryByID = `
		SELECT * FROM delivery
		WHERE delivery_id = $1;	
	`

	getPaymentsByID = `
		SELECT * FROM payments
		WHERE payment_id = $1;
	`

	getItemsByOrderID = `
		SELECT 
			order_id, 
			chrt_id, 
			track_number, 
			price, 
			rid, 
			name, 
			sale, 
			size, 
			total_price, 
			nm_id, 
			brand, 
			status 
		FROM items 
		WHERE order_id = $1;
	`
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *entity.Order) error
	GetAllOrders(ctx context.Context) (map[string]*entity.Order, error)
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

func (r *orderRepository) GetAllOrders(ctx context.Context) (map[string]*entity.Order, error) {
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
			return nil, fmt.Errorf("GetAllOrders: GetDelivery: %s", err.Error())
		}
		if err := r.db.GetContext(ctx, payment, getPaymentsByID, o.PaymentID); err != nil {
			return nil, fmt.Errorf("GetAllOrders: GetPayments: %s", err.Error())
		}
		if err := r.db.SelectContext(ctx, &items, getItemsByOrderID, o.OrderUID); err != nil {
			return nil, fmt.Errorf("GetAllOrders: GetItems: %s", err.Error())
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
