package app

import (
	"context"
	"fmt"
	"wb_tech_l0/internal/infrastructure/broker"
	"wb_tech_l0/internal/infrastructure/cache"
	"wb_tech_l0/internal/infrastructure/config"
	"wb_tech_l0/internal/infrastructure/database"
	"wb_tech_l0/internal/infrastructure/logger"
	"wb_tech_l0/internal/repository"
)

func (a *App) Init() error {
	if err := a.initConfig(); err != nil {
		return err
	}
	if err := a.initLogger(); err != nil {
		return err
	}
	if err := a.initOrderRepository(); err != nil {
		return err
	}
	if err := a.initAndRestoreCache(); err != nil {
		return err
	}
	if err := a.initBroker(); err != nil {
		return err
	}
	return nil
}

func (a *App) initConfig() error {
	configPath := a.getConfigPath()
	if err := config.InitConfig(configPath); err != nil {
		return err
	}
	return nil
}

func (a *App) initLogger() error {
	err := logger.InitLogger(config.C.LoggerConfig.Path)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initBroker() error {
	return broker.InitKafka()
}

func (a *App) initOrderRepository() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.C.PostgresConfig.Host,
		config.C.PostgresConfig.User,
		config.C.PostgresConfig.Password,
		config.C.PostgresConfig.DBName,
		config.C.PostgresConfig.Port,
		config.C.PostgresConfig.SSLmode,
	)

	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		return err
	}

	a.orderRepository = repository.NewOrderRepository(db)
	return nil
}

func (a *App) initAndRestoreCache() error {
	a.memoryCache = cache.NewMemoryCache(config.C.MemoryCacheConfig.Capacity)

	orders, err := a.orderRepository.GetAll(context.Background())
	if err != nil {
		return err
	}

	for _, order := range orders {
		a.memoryCache.Set(order.Key(), order)
	}

	return nil
}
