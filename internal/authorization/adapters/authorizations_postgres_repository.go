package adapters

import (
	"database/sql"
	"errors"
)

type AuthorizationsPostgresRepository struct {
	session *sql.DB
	schema  string
}

func NewAuthorizationsPostgresRepository(session *sql.DB, schema string) (*AuthorizationsPostgresRepository, error) {
	if session == nil {
		return nil, errors.New("db session is nil")
	}

	return &AuthorizationsPostgresRepository{
		session: session,
		schema:  schema,
	}, nil
}

func (r *AuthorizationsPostgresRepository) AddAuthorization(userUUID, token string) error {
	return nil
}

func (r *AuthorizationsPostgresRepository) GetUserUUIDByToken(token string) (string, error) {
	return "", nil
}

func (r *AuthorizationsPostgresRepository) RemoveAuthorization(userUUID, token string) error {
	return nil
}
