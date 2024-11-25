package app

import (
	"context"
	"go.uber.org/zap"
)

func (a *App) shutdown() {
	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		zap.L().Error("Failed http-server shutdown", zap.Error(err))
	}
	if err := a.orderSpamHandler.Shutdown(); err != nil {
		zap.L().Error("Failed order spam handler shutdown", zap.Error(err))
	}
	if err := a.orderProcessHandler.Shutdown(); err != nil {
		zap.L().Error("Failed order process handler shutdown", zap.Error(err))
	}
}
