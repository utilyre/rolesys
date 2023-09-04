package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/utilyre/role/config"
	"github.com/utilyre/role/storage"
)

type Claims struct {
	jwt.RegisteredClaims

	Email string       `json:"email"`
	Role  storage.Role `json:"role"`
}

type Auth struct {
	config config.Config
}

func New(c config.Config) Auth {
	return Auth{config: c}
}

func (a Auth) GenerateToken(email string, role storage.Role) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.config.JWTExpirationTime)),
		},
		Email: email,
		Role:  role,
	}

	token := jwt.NewWithClaims(a.config.JWTSigningMethod, claims)
	return token.SignedString(a.config.JWTSecret)
}

func (a Auth) Allow(roles []storage.Role) (echo.MiddlewareFunc, echo.MiddlewareFunc) {
	jwtware := echojwt.WithConfig(echojwt.Config{
		SigningKey:    a.config.JWTSecret,
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return new(Claims) },
	})

	roleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*Claims)

			for _, role := range roles {
				if role == claims.Role {
					return next(c)
				}
			}

			return echo.ErrForbidden
		}
	}

	return jwtware, roleware
}
