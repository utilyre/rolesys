package main

import (
	"github.com/utilyre/jwtrole/auth"
	"github.com/utilyre/jwtrole/config"
	"github.com/utilyre/jwtrole/database"
	"github.com/utilyre/jwtrole/handler"
	"github.com/utilyre/jwtrole/logger"
	"github.com/utilyre/jwtrole/router"
	"github.com/utilyre/jwtrole/storage"
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
