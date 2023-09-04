package main

import (
	"github.com/utilyre/role/config"
	"github.com/utilyre/role/handler"
	"github.com/utilyre/role/logger"
	"github.com/utilyre/role/router"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.WithLogger(func(l *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: l}
		}),
		fx.Provide(
			config.New,
			logger.New,
			router.New,
		),
		fx.Invoke(
			handler.SetupHealthCheck,
		),
	).Run()
}
