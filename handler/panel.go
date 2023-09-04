package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type panel struct{}

func Panel(e *echo.Echo) {
	g := e.Group("/api/panel")
	p := panel{}

	g.GET("/public", p.panelPublic)
	g.GET("/private", p.panelPrivate)
}

func (p panel) panelPublic(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (p panel) panelPrivate(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
