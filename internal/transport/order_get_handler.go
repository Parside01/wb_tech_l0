package transport

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"wb_tech_l0/internal/service"
)

type OrderGetHandler struct {
	order_service service.OrderService
}

func NewOrderGetHandler(service service.OrderService) *OrderGetHandler {
	return &OrderGetHandler{
		order_service: service,
	}
}

func (c *OrderGetHandler) GetOrder(e echo.Context) error {
	id := e.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing id")
	}
	order, err := c.order_service.GetOrderById(e.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, order)
}
