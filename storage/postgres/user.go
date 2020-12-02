package postgres

import (
	"database/sql"
	"fmt"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/account"
)

// CreateUser inserts a new user with hashed password
func (s *PostgresStorage) CreateUser(id, email, password string) *internal.Error {
	query := fmt.Sprintf(`
		INSERT INTO %s."user"(id, email, password)
		VALUES ($1, $2, $3)
	`, s.schema)

	_, err := s.session.Exec(query, id, email, password)

	if err != nil {
		return internal.NewError(internal.UnknownErrorCode).WithCause(err)
	}

	return nil
}

// GetByEmail find the user with his email
func (s *PostgresStorage) GetByEmail(email string) (*account.User, *internal.Error) {
	query := fmt.Sprintf(`
		SELECT id, password
		FROM %s."user"
		WHERE LOWER("user".email) = LOWER($1)
	`, s.schema)

	row := s.session.QueryRow(query, email)

	u := &account.User{}

	err := row.Scan(&u.ID, &u.Password)

	if err == sql.ErrNoRows {
		return nil, internal.NewError(account.UserNotFoundErrorCode)
	}

	if err != nil {
		return nil, internal.NewError(internal.UnknownErrorCode).WithCause(err)
	}

	return u, nil
}

// GetUserIDByToken find the user id with his token
func (s *PostgresStorage) GetUserIDByToken(token string) (*string, *internal.Error) {
	query := fmt.Sprintf(`
		SELECT id
    FROM %s."user"
    LEFT JOIN %s."authorization" ON "authorization".user_id = "user".id
		WHERE "authorization".token = $1
	`, s.schema, s.schema)

	row := s.session.QueryRow(query, token)

	var id *string

	err := row.Scan(&id)

	if err == sql.ErrNoRows {
		return nil, internal.NewError(account.UserNotFoundErrorCode)
	}

	if err != nil {
		return nil, internal.NewError(internal.UnknownErrorCode).WithCause(err)
	}

	return id, nil
}

// GetByToken find the user with his token
func (s *PostgresStorage) GetByToken(token string) (*account.User, *internal.Error) {
	query := fmt.Sprintf(`
		SELECT id, password, email
    FROM %s."user"
    LEFT JOIN %s."authorization" ON "authorization".user_id = "user".id
		WHERE "authorization".token = $1
	`, s.schema, s.schema)

	row := s.session.QueryRow(query, token)

	u := &account.User{}

	err := row.Scan(&u.ID, &u.Password, &u.Email)

	if err == sql.ErrNoRows {
		return nil, internal.NewError(account.UserNotFoundErrorCode)
	}

	if err != nil {
		return nil, internal.NewError(internal.UnknownErrorCode).WithCause(err)
	}

	return u, nil
}
