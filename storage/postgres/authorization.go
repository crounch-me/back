package postgres

import (
	"fmt"

	"github.com/crounch-me/back/domain"
)

// CreateAuthorization creates a user id and token couple
func (s *PostgresStorage) CreateAuthorization(userID, token string) *domain.Error {
	query := fmt.Sprintf(`
		INSERT INTO %s."authorization" (user_id, token)
		VALUES ($1, $2)
	`, s.schema)

	_, err := s.session.Exec(query, userID, token)

	if err != nil {
		return domain.NewError(domain.UnknownErrorCode).WithCause(err)
	}

	return nil
}

func (s *PostgresStorage) DeleteAuthorization(userID, token string) *domain.Error {
  query := fmt.Sprintf(`
    DELETE FROM %s.authorization WHERE user_id = $1 AND token = $2
  `, s.schema)

  _, err := s.session.Exec(query, userID, token)

  if err != nil {
    return domain.NewError(domain.UnknownErrorCode).WithCause(err)
  }

  return nil
}
