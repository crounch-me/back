package authorization

import (
	"github.com/crounch-me/back/domain"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (sm *StorageMock) CreateAuthorization(userID, token string) *domain.Error {
	args := sm.Called(userID, token)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*domain.Error)
}

func (sm *StorageMock) DeleteAuthorization(userID, token string) *domain.Error {
	args := sm.Called(userID, token)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*domain.Error)
}
