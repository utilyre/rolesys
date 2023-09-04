package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func SetupHealthCheck(l *zap.Logger, e *echo.Echo) {
	healthcheck := e.Group("/api/healthcheck")

	healthcheck.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
}
