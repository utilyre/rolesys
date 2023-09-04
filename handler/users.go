package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utilyre/jwtrole/auth"
	"github.com/utilyre/jwtrole/storage"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint64 `json:"id" validate:"isdefault"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=1024"`
}

type users struct {
	auth    auth.Auth
	storage storage.UsersStorage
}

func Users(e *echo.Echo, a auth.Auth, s storage.UsersStorage) {
	g := e.Group("/api/users")
	u := users{auth: a, storage: s}

	g.POST("/signup", u.usersSignUp)
	g.POST("/signin", u.usersSignIn)
}

func (u users) usersSignUp(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	if err := c.Validate(user); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	sUser := &storage.User{
		Email:    user.Email,
		Password: hash,
		Role:     storage.RoleUser,
	}
	if err := u.storage.Create(sUser); err != nil {
		if errors.Is(err, storage.ErrDuplicateKey) {
			return echo.ErrConflict
		}

		return err
	}

	user.ID = sUser.ID
	return c.JSON(http.StatusCreated, user)
}

func (u users) usersSignIn(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	if err := c.Validate(user); err != nil {
		return err
	}

	sUser, err := u.storage.GetByEmail(user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.ErrNotFound
		}

		return err
	}

	if err := bcrypt.CompareHashAndPassword(sUser.Password, []byte(user.Password)); err != nil {
		return echo.ErrNotFound
	}

	token, err := u.auth.GenerateToken(sUser.Email, sUser.Role)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{"token": token})
}
