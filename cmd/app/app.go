package app

import (
	"wb_tech_l0/internal/infrastructure/cache"
	"wb_tech_l0/internal/infrastructure/database"
	"wb_tech_l0/internal/infrastructure/logger"
	"wb_tech_l0/internal/repository"
	"wb_tech_l0/internal/service"
)

type App struct{}

func (a *App) Start() error {
	l, err := logger.NewLogger("log/log.log")
	if err != nil {
		return err
	}

	repo, err := a.setupDBConnection()
	c := cache.NewMemoryCache()
	serv := service.NewOrderService(repo)
}

func (a *App) setupDBConnection() (repository.OrderRepository, error) {
	db, err := database.NewPostgresDB()
	if err != nil {
		return nil, err
	}
	repo := repository.NewOrderRepository(db)
	return repo, nil
}
