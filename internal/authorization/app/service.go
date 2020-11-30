package app

import (
	"errors"

	"github.com/crounch-me/back/internal/authorization/domain/authorizations"
)

type AuthorizationService struct {
	authorizationsRepository authorizations.Repository
}

func NewAuthorizationService(authorizationsRepository authorizations.Repository) (*AuthorizationService, error) {
	if authorizationsRepository == nil {
		return nil, errors.New("authorizationsRepository is nil")
	}

	return &AuthorizationService{
		authorizationsRepository: authorizationsRepository,
	}, nil
}

func (s *AuthorizationService) GetUserUUIDByToken(token string) (string, error) {
	return s.authorizationsRepository.GetUserUUIDByToken(token)
}
