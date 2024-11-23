package transport

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

func LoggingMiddlewareEcho(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		start := time.Now()
		err := next(e)
		duration := time.Since(start)

		zap.L().Info("Handled request",
			zap.String("method", e.Request().Method),
			zap.String("path", e.Request().URL.Path),
			zap.String("duration", duration.String()),
		)

		if err != nil {
			zap.L().Error("Request failed",
				zap.Error(err),
			)
		}
		return err
	}
}
