package authorization

import "github.com/crounch-me/back/domain"

type Storage interface {
	CreateAuthorization(userID, token string) *domain.Error
}
