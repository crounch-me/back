package postgres

import (
	"database/sql"
	"fmt"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
)

// CreateUser inserts a new user with hashed password
func (s *PostgresStorage) CreateUser(id, email, password string) *domain.Error {
	query := fmt.Sprintf(`
		INSERT INTO %s."user"(id, email, password)
		VALUES ($1, $2, $3)
	`, s.schema)

	_, err := s.session.Exec(query, id, email, password)

	if err != nil {
		return domain.NewError(domain.UnknownErrorCode).WithCause(err)
	}

	return nil
}

// GetByEmail find the user with his email
func (s *PostgresStorage) GetByEmail(email string) (*users.User, *domain.Error) {
	query := fmt.Sprintf(`
		SELECT id, password
		FROM %s."user"
		WHERE LOWER("user".email) = LOWER($1)
	`, s.schema)

	row := s.session.QueryRow(query, email)

	user := &users.User{}

	err := row.Scan(&user.ID, &user.Password)

	if err == sql.ErrNoRows {
		return nil, domain.NewError(users.UserNotFoundErrorCode)
	}

	if err != nil {
		return nil, domain.NewError(domain.UnknownErrorCode).WithCause(err)
	}

	return user, nil
}

// GetUserIDByToken find the user id with his token
func (s *PostgresStorage) GetUserIDByToken(token string) (*string, *domain.Error) {
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
		return nil, domain.NewError(users.UserNotFoundErrorCode)
	}

	if err != nil {
		return nil, domain.NewError(domain.UnknownErrorCode).WithCause(err)
	}

	return id, nil
}

// GetByToken find the user with his token
func (s *PostgresStorage) GetByToken(token string) (*users.User, *domain.Error) {
	query := fmt.Sprintf(`
		SELECT id, email
    FROM %s."user"
    LEFT JOIN %s."authorization" ON "authorization".user_id = "user".id
		WHERE "authorization".token = $1
	`, s.schema, s.schema)

	row := s.session.QueryRow(query, token)

	user := &users.User{}

	err := row.Scan(&user.ID, &user.Password)

	if err == sql.ErrNoRows {
		return nil, domain.NewError(users.UserNotFoundErrorCode)
	}

	if err != nil {
		return nil, domain.NewError(domain.UnknownErrorCode).WithCause(err)
	}

	return user, nil
}
