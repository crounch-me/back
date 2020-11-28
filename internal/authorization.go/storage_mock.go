package authorization

import (
	"github.com/crounch-me/back/internal"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (sm *StorageMock) CreateAuthorization(userID, token string) *internal.Error {
	args := sm.Called(userID, token)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*internal.Error)
}

func (sm *StorageMock) DeleteAuthorization(userID, token string) *internal.Error {
	args := sm.Called(userID, token)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*internal.Error)
}
