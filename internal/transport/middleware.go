package transport

import (
	"github.com/labstack/echo/v4"
	"time"
	"wb_tech_l0/internal/transport/metrics"
)

func RequestDurationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(c.Request().URL.Path).Observe(duration)
		return err
	}
}

func RequestCountMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		metrics.HttpRequestCountWithPath.WithLabelValues(c.Request().URL.Path).Inc()
		return next(c)
	}
}
