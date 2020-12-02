package adapters

import (
	"database/sql"
	"errors"

	"github.com/crounch-me/back/internal/account/domain/users"
)

type UsersPostgresRepository struct {
	session *sql.DB
	schema  string
}

func NewUsersPostgresRepository(session *sql.DB, schema string) (*UsersPostgresRepository, error) {
	if session == nil {
		return nil, errors.New("db session is nil")
	}

	return &UsersPostgresRepository{
		session: session,
		schema:  schema,
	}, nil
}

func (r *UsersPostgresRepository) AddUser(user *users.User) error {
	return nil
}

func (r *UsersPostgresRepository) FindByEmail(email string) (*users.User, error) {
	return nil, nil
}
