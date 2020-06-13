package postgres

import (
	"fmt"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/authorization.go"
	"github.com/crounch-me/back/domain/users"
)

// CreateAuthorization creates a user id and token couple
func (s *PostgresStorage) CreateAuthorization(userID, token string) (*authorization.Authorization, *domain.Error) {
	query := fmt.Sprintf(`
		INSERT INTO %s."authorization" (user_id, token)
		VALUES ($1, $2)
	`, s.schema)

	_, err := s.session.Exec(query, userID, token)

	if err != nil {
		return nil, domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}

	authorization := &authorization.Authorization{
		AccessToken: token,
		Owner: &users.User{
			ID: userID,
		},
	}

	return authorization, nil
}
