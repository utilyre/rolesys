package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/utilyre/rolesys/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Validator struct {
	validate validator.Validate
}

var _ echo.Validator = (*Validator)(nil)

func (v *Validator) Validate(i any) error {
	if err := v.validate.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return nil
}

func New(lc fx.Lifecycle, c config.Config, l *zap.Logger) *echo.Echo {
	e := echo.New()

	e.Debug = true
	e.HideBanner = true
	e.Validator = &Validator{validate: *validator.New()}

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
