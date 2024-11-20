package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"wb_tech_l0/internal/entity"
	"wb_tech_l0/internal/infrastructure/database"
	"wb_tech_l0/internal/repository"
	service "wb_tech_l0/internal/service"
)

func main() {
	db, err := database.NewPostgresDB()
	if err != nil {
		panic(err)
	}
	repo := repository.NewOrderRepository(db)
	service := service.NewOrderService(repo)

	router := echo.New()
	router.POST("/", func(c echo.Context) error {
		var order *entity.Order
		if err := c.Bind(&order); err != nil {
			return err
		}
		service.SaveOrder(c.Request().Context(), order)
		return nil
	})

	router.GET("/", func(c echo.Context) error {
		all, err := repo.GetAllOrders(c.Request().Context())
		if err != nil {
			fmt.Println(err)
			return err
		}
		return c.JSON(http.StatusOK, all)
	})
	if err := router.Start(":8080"); err != nil {
		panic(err)
	}
}
