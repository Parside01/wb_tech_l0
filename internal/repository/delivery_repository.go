package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"wb_tech_l0/internal/entity"
)

type DeliveryRepository interface {
	SaveTX(ctx context.Context, tx *sqlx.Tx, delivery *entity.Delivery) error
	GetByID(ctx context.Context, id string) (*entity.Delivery, error)
}

type deliveryRepository struct {
	db *sqlx.DB
}

func NewDeliveryRepository(db *sqlx.DB) DeliveryRepository {
	return &deliveryRepository{
		db: db,
	}
}

func (d *deliveryRepository) SaveTX(ctx context.Context, tx *sqlx.Tx, delivery *entity.Delivery) error {
	if _, err := tx.ExecContext(ctx, saveDelivery,
		delivery.DeliveryID,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	); err != nil {
		return err
	}
	return nil
}

func (d *deliveryRepository) GetByID(ctx context.Context, id string) (*entity.Delivery, error) {
	delivery := &entity.Delivery{}
	if err := d.db.GetContext(ctx, delivery, getDeliveryByID, id); err != nil {
		return nil, err
	}
	return delivery, nil
}
