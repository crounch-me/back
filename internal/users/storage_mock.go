package users

import (
	"github.com/crounch-me/back/internal"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (sm *StorageMock) CreateUser(id, email, password string) *internal.Error {
	args := sm.Called(id, email, password)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*internal.Error)
}

func (sm *StorageMock) GetUserIDByToken(token string) (*string, *internal.Error) {
	args := sm.Called(token)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*string), nil
	} else {
		return nil, err.(*internal.Error)
	}
}

func (sm *StorageMock) GetByEmail(email string) (*User, *internal.Error) {
	args := sm.Called(email)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*User), nil
	}
	return nil, err.(*internal.Error)
}

func (sm *StorageMock) GetByToken(token string) (*User, *internal.Error) {
	args := sm.Called(token)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*User), nil
	}
	return nil, err.(*internal.Error)
}
