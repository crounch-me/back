package authorization

import "github.com/crounch-me/back/internal"

type Storage interface {
	CreateAuthorization(userID, token string) *internal.Error
	DeleteAuthorization(userID, token string) *internal.Error
}
