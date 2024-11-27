package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"wb_tech_l0/internal/entity"
)

func TestDeliveryRepository_SaveTX(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewDeliveryRepository(db)

	delivery := entity.GenerateRandomDelivery()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO delivery").
		WithArgs(
			delivery.OrderID,
			delivery.Name,
			delivery.Phone,
			delivery.Zip,
			delivery.City,
			delivery.Address,
			delivery.Region,
			delivery.Email).
		WillReturnResult(sqlxmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}
	err = repo.SaveTX(context.Background(), tx, delivery)

	assert.NoError(t, err)
}

func TestDeliveryRepository_GetByOrderID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewDeliveryRepository(db)
	delivery := entity.GenerateRandomDelivery()

	mock.ExpectQuery("SELECT (.+) FROM delivery WHERE order_id = ?").
		WithArgs(delivery.OrderID).
		WillReturnRows(sqlxmock.NewRows([]string{
			"order_id", "name", "phone", "zip", "city", "address", "region", "email",
		}).AddRow(delivery.OrderID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email))

	got, err := repo.GetByOrderID(context.Background(), delivery.OrderID)

	assert.NoError(t, err)
	assert.Equal(t, delivery, got)
}
