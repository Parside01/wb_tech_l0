package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
	"sync"
	"wb_tech_l0/internal/infrastructure/cache"
	"wb_tech_l0/internal/infrastructure/config"
	"wb_tech_l0/internal/repository"
	"wb_tech_l0/internal/service"
	"wb_tech_l0/internal/transport"
)

type App struct {
	orderRepository     repository.OrderRepository
	memoryCache         cache.Cache
	orderService        service.OrderService
	orderProcessHandler *transport.OrderProcessHandler
	orderSpamHandler    *transport.OrderSpamHandler
	orderGetHandler     *transport.OrderGetHandler
	httpServer          *echo.Echo
	errorsChan          chan error
	wg                  sync.WaitGroup
}

func (a *App) Start() error {
	if err := a.setupHandlers(); err != nil {
		return err
	}
	a.setupHttpServer()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		address := fmt.Sprintf("%s:%s", config.C.HttpServerConfig.Host, config.C.HttpServerConfig.Port)
		if err := a.httpServer.Start(address); err != nil {
			return
		}
	}()
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		if err := a.orderProcessHandler.Start(context.Background()); err != nil {
			return
		}
	}()
	a.wg.Wait()
	return nil
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
