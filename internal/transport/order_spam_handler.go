package transport

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
	"wb_tech_l0/internal/entity"
	"wb_tech_l0/internal/infrastructure/broker"
)

type OrderSpamHandler struct {
	publisher *broker.KafkaPublisher
}

func NewOrderSpamHandler(publisher *broker.KafkaPublisher) *OrderSpamHandler {
	return &OrderSpamHandler{
		publisher: publisher,
	}
}

func (c *OrderSpamHandler) SpamOrders(e echo.Context) error {
	count64, err := strconv.ParseInt(e.QueryParam("count"), 10, 32)
	if err != nil {
		e.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	count := int(count64)

	for i := 0; i < count; i++ {
		order := entity.GenerateRandomOrder()
		data, err := json.Marshal(order)
		if err != nil {
			e.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		key, err := json.Marshal(order.Key())
		if err != nil {
			e.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		start := time.Now()
		if err := c.publisher.PublishMessage(e.Request().Context(), key, data); err != nil {
			e.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		zap.L().Warn("Success publish order", zap.String("ID", order.OrderUID), zap.String("duration", time.Since(start).String()))
	}
	return e.JSON(http.StatusOK, fmt.Sprintf("Success publish %d orders", count))
}
