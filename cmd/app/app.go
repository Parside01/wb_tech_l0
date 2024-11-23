package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"os"
	"wb_tech_l0/internal/infrastructure/broker"
	"wb_tech_l0/internal/infrastructure/cache"
	"wb_tech_l0/internal/infrastructure/config"
	"wb_tech_l0/internal/infrastructure/database"
	"wb_tech_l0/internal/infrastructure/logger"
	"wb_tech_l0/internal/repository"
	"wb_tech_l0/internal/service"
	"wb_tech_l0/internal/transport"
)

type App struct {
	log           *zap.Logger
	order_handler *transport.OrderProcessHandler
	server        *echo.Echo
}

func (a *App) Init() error {
	configPath := a.getConfigPath()
	if err := config.InitConfig(configPath); err != nil {
		return err
	}
	if err := broker.InitKafka(); err != nil {
		return err
	}
	return nil
}

func (a *App) Start() error {
	l, err := logger.NewLogger("./log/log.log")
	if err != nil {
		return err
	}
	a.log = l
	l.Info("Logger initialization successful")
	if err := a.setupOrderProcessHandler(); err != nil {
		return err
	}

	a.setupHttpServer()
	if err := a.server.Start(fmt.Sprintf("%s:%s", config.Global.HttpServerConfig.Host, config.Global.HttpServerConfig.Port)); err != nil {
		return err
	}

	if err := a.order_handler.Start(context.Background()); err != nil {
		return err
	}
	return nil
}

func (a *App) setupHttpServer() {
	a.server = echo.New()
	publisher := broker.NewKafkaPublisher()

	spamHandler := transport.NewOrderSpamHandler(a.server, publisher)
	a.server.POST("/order/spam", spamHandler.SpamOrders, middleware.Logger())
}

func (a *App) setupOrderProcessHandler() error {
	repo, err := a.setupDBConnection()
	if err != nil {
		return err
	}
	c := cache.NewMemoryCache(config.Global.MemoryCacheConfig.Capacity)

	if err := a.restoreCache(repo, c); err != nil {
		return err
	}

	serv := service.NewOrderService(repo, c)
	cons := broker.NewKafkaConsumer()
	a.order_handler = transport.NewOrderProcessHandler(serv, cons, a.log)
	return nil
}

func (a *App) restoreCache(repo repository.OrderRepository, c cache.Cache) error {
	orders, err := repo.GetAll(context.Background())
	if err != nil {
		return err
	}

	for _, order := range orders {
		c.Set(order.Key(), order)
	}
	return nil
}

func (a *App) setupDBConnection() (repository.OrderRepository, error) {
	fmt.Println(*config.Global.PostgresConfig)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Global.PostgresConfig.Host,
		config.Global.PostgresConfig.User,
		config.Global.PostgresConfig.Password,
		config.Global.PostgresConfig.DBName,
		config.Global.PostgresConfig.Port,
		config.Global.PostgresConfig.SSLmode,
	)

	db, err := database.NewPostgresDB(dsn)
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
