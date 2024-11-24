package transport

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
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

	messages := make([]broker.KafkaMessage, count64)
	response := make([]string, count64)
	for i := range messages {
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
		messages[i] = broker.KafkaMessage{
			Key:   key,
			Value: data,
		}
		response[i] = order.Key()
	}
	if err := c.publisher.PublishMessages(e.Request().Context(), messages...); err != nil {
		zap.L().Error("Failed publish messages", zap.Int64("Message count", count64), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return e.JSON(http.StatusOK, response)
}
