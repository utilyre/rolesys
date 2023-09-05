package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/utilyre/rolesys/config"
	"github.com/utilyre/rolesys/internal/cookies"
	"github.com/utilyre/rolesys/storage"
)

const cookieName = "Session"

type Auth struct {
	config  config.Config
	storage storage.SessionsStorage
}

func New(c config.Config, s storage.SessionsStorage) Auth {
	return Auth{config: c, storage: s}
}

func (a Auth) WriteToken(w http.ResponseWriter, user_id uint64) error {
	token := uuid.NewString()
	expiresAt := time.Now().Add(a.config.AuthExpirationTime)

	if err := a.storage.Create(&storage.Session{
		UserID:    user_id,
		Token:     token,
		ExpiresAt: &expiresAt,
	}); err != nil {
		return err
	}

	return cookies.WriteEncrypted(
		w,
		&http.Cookie{
			Name:     cookieName,
			Value:    token,
			Expires:  expiresAt,
			Path:     "/api",
			HttpOnly: true,
			Secure:   true,
		},
		a.config.AuthSecret,
	)
}

func (a Auth) Allow(roles []storage.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := cookies.ReadEncrypted(c.Request(), cookieName, a.config.AuthSecret)
			if err != nil {
				return err
			}

			session, err := a.storage.GetJoinedUsersByToken(cookie.Value)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return echo.ErrNotFound
				}

				return err
			}

			if session.ExpiresAt.Before(time.Now()) {
				return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
			}

			c.Set("user", &session.User)
			return next(c)
		}
	}
}
