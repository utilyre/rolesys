package main

import (
	"github.com/utilyre/role/auth"
	"github.com/utilyre/role/config"
	"github.com/utilyre/role/database"
	"github.com/utilyre/role/handler"
	"github.com/utilyre/role/logger"
	"github.com/utilyre/role/router"
	"github.com/utilyre/role/storage"
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
			database.New,
			auth.New,
			router.New,

			storage.NewUsers,
		),
		fx.Invoke(
			handler.Users,
			handler.Panel,
		),
	).Run()
}
