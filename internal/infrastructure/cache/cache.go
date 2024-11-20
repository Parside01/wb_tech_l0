package cache

import (
	"wb_tech_l0/internal/entity"
)

type Cache interface {
	Set(key string, order *entity.Order)
	Get(key string) (*entity.Order, bool)
	Delete(key string)
}
