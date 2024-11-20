package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"wb_tech_l0/internal/entity"
)

type ItemRepository interface {
	SaveTX(ctx context.Context, tx *sqlx.Tx, item *entity.Item) error
	GetAllByOrderID(ctx context.Context, orderID string) ([]*entity.Item, error)
}

type itemRepository struct {
	db *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) ItemRepository {
	return &itemRepository{
		db: db,
	}
}

func (r *itemRepository) SaveTX(ctx context.Context, tx *sqlx.Tx, item *entity.Item) error {
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
		return err
	}
	return nil
}

func (r *itemRepository) GetAllByOrderID(ctx context.Context, orderID string) ([]*entity.Item, error) {
	items := []*entity.Item{}

	if err := r.db.SelectContext(ctx, &items, getItemsByOrderID, orderID); err != nil {
		return nil, err
	}
	return items, nil
}
