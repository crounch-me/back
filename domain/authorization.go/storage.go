package authorization

import "github.com/crounch-me/back/domain"

type AuthorizationStorage interface {
	CreateAuthorization(userID, token string) (*Authorization, *domain.Error)
}
