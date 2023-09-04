package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/utilyre/role/config"
	"github.com/utilyre/role/storage"
)

type Auth struct {
	config config.Config
}

func New(c config.Config) Auth {
	return Auth{config: c}
}

func (a Auth) GenerateToken(email string, role storage.Role) (string, error) {
	return jwt.NewWithClaims(
		a.config.JWTSigningMethod,
		jwt.MapClaims{
			"email": email,
			"role":  role,
			"exp":   time.Now().Add(a.config.JWTExpirationTime).Unix(),
		},
	).SignedString(a.config.JWTSecret)
}
