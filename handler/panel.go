package handler

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/utilyre/jwtrole/auth"
	"github.com/utilyre/jwtrole/config"
	"github.com/utilyre/jwtrole/storage"
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
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.Claims)

	return c.String(http.StatusOK, claims.Email)
}
