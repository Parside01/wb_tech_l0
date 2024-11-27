package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"wb_tech_l0/internal/entity"
)

func TestItemRepository_SaveTX(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewItemRepository(db)
	item := entity.GenerateRandomItems()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO items").WithArgs(
		item[0].OrderID,
		item[0].ChrtID,
		item[0].TrackNumber,
		item[0].Price,
		item[0].RID,
		item[0].Name,
		item[0].Sale,
		item[0].Size,
		item[0].TotalPrice,
		item[0].NMID,
		item[0].Brand,
		item[0].Status).WillReturnResult(sqlxmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}
	err = repo.SaveTX(context.Background(), tx, item[0])

	assert.NoError(t, err)
}

func TestItemRepository_GetAllByOrderID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewItemRepository(db)
	item := entity.GenerateRandomItem()

	mock.ExpectQuery("SELECT (.+) FROM items WHERE order_id = ?").
		WithArgs(item.OrderID).
		WillReturnRows(sqlxmock.NewRows([]string{"order_id", "chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status"}).
			AddRow(item.OrderID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status))

	got, err := repo.GetAllByOrderID(context.Background(), item.OrderID)

	assert.NoError(t, err)
	assert.Equal(t, []*entity.Item{item}, got)
}
