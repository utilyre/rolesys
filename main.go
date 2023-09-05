package main

import (
	"github.com/utilyre/rolesys/auth"
	"github.com/utilyre/rolesys/config"
	"github.com/utilyre/rolesys/database"
	"github.com/utilyre/rolesys/handler"
	"github.com/utilyre/rolesys/logger"
	"github.com/utilyre/rolesys/router"
	"github.com/utilyre/rolesys/storage"
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
			storage.NewSessions,
		),
		fx.Invoke(
			handler.Users,
			handler.Panel,
		),
	).Run()
}
