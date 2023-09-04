package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utilyre/role/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(lc fx.Lifecycle, c config.Config, l *zap.Logger) *echo.Echo {
	e := echo.New()

	e.HideBanner = true

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := e.Start(fmt.Sprintf("%s:%s", c.BEHost, c.BEPort)); err != nil {
					switch {
					case errors.Is(err, http.ErrServerClosed):
						l.Info("HTTP server closed")
					default:
						l.Error("failed to start HTTP server", zap.Error(err))
					}
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	return e
}
