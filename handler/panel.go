package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utilyre/rolesys/auth"
	"github.com/utilyre/rolesys/config"
	"github.com/utilyre/rolesys/storage"
)

type panel struct{}

func Panel(e *echo.Echo, c config.Config, a auth.Auth) {
	g := e.Group("/api/panel")
	p := panel{}

	g.GET("/public", p.panelPublic)
	g.Use(a.Allow([]storage.Role{storage.RoleUser, storage.RoleAdmin}))
	g.GET("/private", p.panelPrivate)
}

func (p panel) panelPublic(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (p panel) panelPrivate(c echo.Context) error {
	user := c.Get("user").(*storage.User)

	return c.String(http.StatusOK, user.Email)
}
