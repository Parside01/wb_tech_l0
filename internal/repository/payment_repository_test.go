package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"wb_tech_l0/internal/entity"
)

func TestPaymentRepository_SaveTX(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPaymentRepository(db)

	payment := entity.GenerateRandomPayment()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO payments").
		WithArgs(
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
		).
		WillReturnResult(sqlxmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}
	err = repo.SaveTX(context.Background(), tx, payment)

	assert.NoError(t, err)
}

func TestPaymentRepository_GetByOrderID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewPaymentRepository(db)
	payment := entity.GenerateRandomPayment()

	mock.ExpectQuery("SELECT (.+) FROM payments WHERE order_id = ?").
		WithArgs(payment.OrderID).
		WillReturnRows(sqlxmock.NewRows([]string{"order_id", "transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee"}).
			AddRow(payment.OrderID, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDT, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee))

	got, err := repo.GetByOrderID(context.Background(), payment.OrderID)

	assert.NoError(t, err)
	assert.Equal(t, payment, got)
}
