package storage

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrDuplicateKey = errors.New("duplicate key value violates unique constraint")
)

type Role int

const (
	RoleGuest Role = iota
	RoleUser
	RoleAdmin
)

type User struct {
	ID       uint64 `db:"id"`
	Email    string `db:"email"`
	Password []byte `db:"password"`
	Role     Role   `db:"role"`
}

type UsersStorage struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) UsersStorage {
	return UsersStorage{db: db}
}

func (s UsersStorage) Create(user *User) error {
	query := `
	INSERT INTO "users"
	("email", "password", "role")
	VALUES ($1, $2, $3)
	RETURNING "id";
	`

	var id uint64
	if err := s.db.QueryRow(query, user.Email, user.Password, user.Role).Scan(&id); err != nil {
		pqErr := new(pq.Error)
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrDuplicateKey
		}

		return err
	}

	user.ID = id
	return nil
}

func (s UsersStorage) GetByEmail(email string) (*User, error) {
	query := `
	SELECT *
	FROM "users"
	WHERE "email" = $1;
	`

	user := new(User)
	if err := s.db.Get(user, query, email); err != nil {
		return nil, err
	}

	return user, nil
}
