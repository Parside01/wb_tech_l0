package app

import (
	"flag"
	"os"
	"wb_tech_l0/internal/infrastructure/config"
	"wb_tech_l0/internal/infrastructure/database"
	"wb_tech_l0/internal/infrastructure/logger"
	"wb_tech_l0/internal/repository"
)

type App struct{}

func (a *App) Init() error {
	configPath := a.getConfigPath()
	if err := config.InitConfig(configPath); err != nil {
		return err
	}

	return nil
}

func (a *App) Start() error {
	l, err := logger.NewLogger("./log/log.log")
	if err != nil {
		return err
	}
	l.Info("Logger initialization successful")
	//repo, err := a.setupDBConnection()
	//c := cache.NewMemoryCache()
	//serv := service.NewOrderService(repo)

	return nil
}

func (a *App) setupDBConnection() (repository.OrderRepository, error) {
	db, err := database.NewPostgresDB()
	if err != nil {
		return nil, err
	}
	repo := repository.NewOrderRepository(db)
	return repo, nil
}

func (a *App) getConfigPath() string {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")
	if configPath != "" {
		return configPath
	}

	flag.StringVar(&configPath, "config", "configs/static-config.yaml", "path to config file")
	return configPath
}
