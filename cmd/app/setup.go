package app

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"wb_tech_l0/internal/infrastructure/broker"
	"wb_tech_l0/internal/service"
	"wb_tech_l0/internal/transport"
)

func (a *App) setupHandlers() error {
	if err := a.setupOrderProcessHandler(); err != nil {
		return err
	}
	if err := a.setupOrderSpamHandler(); err != nil {
		return err
	}
	if err := a.setupOrderGetHandler(); err != nil {
		return err
	}
	zap.L().Info("Finished setup handlers")
	return nil
}

func (a *App) setupHttpServer() {
	a.httpServer = echo.New()

	a.httpServer.GET(
		"/metrics",
		echo.WrapHandler(promhttp.Handler()),
	)

	v1 := a.httpServer.Group("/api/v1/order")
	v1.POST(
		"/spam",
		a.orderSpamHandler.SpamOrders,
		transport.LoggingMiddlewareEcho,
	)

	v1.GET(
		"/:id",
		a.orderGetHandler.GetOrder,
		transport.LoggingMiddlewareEcho,
	)
}

func (a *App) setupOrderGetHandler() error {
	a.orderGetHandler = transport.NewOrderGetHandler(a.orderService)
	return nil
}

func (a *App) setupOrderProcessHandler() error {
	a.orderService = service.NewOrderService(a.orderRepository, a.memoryCache)
	consumer := broker.NewKafkaConsumer()
	a.orderProcessHandler = transport.NewOrderProcessHandler(a.orderService, consumer)

	return nil
}

func (a *App) setupOrderSpamHandler() error {
	publisher := broker.NewKafkaPublisher()
	a.orderSpamHandler = transport.NewOrderSpamHandler(publisher)

	return nil
}
