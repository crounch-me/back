package mock

import (
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/authorization.go"
)

// CreateAuthorization mocks create authorization
func (sm *StorageMock) CreateAuthorization(userID, token string) (*authorization.Authorization, *domain.Error) {
	args := sm.Called(userID, token)
	return args.Get(0).(*authorization.Authorization), args.Error(1).(*domain.Error)
}
