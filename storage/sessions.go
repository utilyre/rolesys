package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Session struct {
	ID uint64 `db:"id"`

	UserID uint64 `db:"user_id"`
	User   User   `db:"user"`

	Token     string     `db:"token"`
	ExpiresAt *time.Time `db:"expires_at"`
}

type SessionsStorage struct {
	db *sqlx.DB
}

func NewSessions(db *sqlx.DB) SessionsStorage {
	return SessionsStorage{db: db}
}

func (s SessionsStorage) Create(session *Session) error {
	query := `
	INSERT INTO "sessions"
	("user_id", "token", "expires_at")
	VALUES ($1, $2, $3)
	RETURNING "id";
	`

	var id uint64
	if err := s.db.QueryRow(query, session.UserID, session.Token, session.ExpiresAt).Scan(&id); err != nil {
		return err
	}

	session.ID = id
	return nil
}

func (s SessionsStorage) GetJoinedUsersByToken(token string) (*Session, error) {
	query := `
	SELECT s."id" AS "id", s."token" AS "token", s."expires_at" AS "expires_at", u."id" AS "user.id", u."email" AS "user.email", u."role" AS "user.role"
	FROM "sessions" AS s
	INNER JOIN "users" AS u
	ON s."user_id" = u."id"
	WHERE s."token" = $1;
	`

	session := new(Session)
	if err := s.db.Get(session, query, token); err != nil {
		return nil, err
	}

	return session, nil
}
