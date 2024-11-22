package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"wb_tech_l0/internal/entity"
)

type PaymentRepository interface {
	SaveTX(ctx context.Context, tx *sqlx.Tx, payment *entity.Payment) error
	GetByOrderID(ctx context.Context, id string) (*entity.Payment, error)
}

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (p *paymentRepository) SaveTX(ctx context.Context, tx *sqlx.Tx, payment *entity.Payment) error {
	if _, err := tx.ExecContext(ctx, savePayment,
		payment.OrderID,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDT,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	); err != nil {
		return err
	}
	return nil
}

func (p *paymentRepository) GetByOrderID(ctx context.Context, id string) (*entity.Payment, error) {
	payment := &entity.Payment{}
	if err := p.db.GetContext(ctx, payment, getPaymentByOrderID, id); err != nil {
		return nil, err
	}
	return payment, nil
}
