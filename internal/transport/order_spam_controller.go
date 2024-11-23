package transport

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"time"
	"wb_tech_l0/internal/entity"
	"wb_tech_l0/internal/infrastructure/broker"
)

type OrderSpamHandler struct {
	publisher *broker.KafkaPublisher
}

func NewOrderSpamHandler(e *echo.Echo, publisher *broker.KafkaPublisher) *OrderSpamHandler {
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
		point1 := time.Now()
		order := entity.GenerateRandomOrder()
		point2 := time.Now()
		data, err := json.Marshal(order)
		if err != nil {
			e.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		point3 := time.Now()
		key, err := json.Marshal(order.Key())
		if err != nil {
			e.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if err := c.publisher.PublishMessage(e.Request().Context(), key, data); err != nil {
			e.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		point4 := time.Now()
		log.Warnf("Generate: %d  Marshal: %d Publish: %d", point2.Unix()-point1.Unix(), point3.Unix()-point2.Unix(), point4.Unix()-point3.Unix())
	}
	return e.JSON(http.StatusOK, fmt.Sprintf("Success publish %d orders", count))
}
