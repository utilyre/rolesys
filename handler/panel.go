package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Panel(e *echo.Echo) {
	panel := e.Group("/api/panel")

	panel.GET("/public", panelPublic())
	panel.GET("/private", panelPrivate())
}

func panelPublic() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}
}

func panelPrivate() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}
}
