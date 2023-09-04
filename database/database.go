package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/utilyre/jwtrole/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(lc fx.Lifecycle, c config.Config, l *zap.Logger) *sqlx.DB {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName,
		),
	)
	if err != nil {
		l.Fatal("failed to establish database connection", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return db.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db
}
