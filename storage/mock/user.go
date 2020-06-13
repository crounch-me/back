package mock

import (
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
)

func (sm *StorageMock) CreateUser(u *users.User) *domain.Error {
	args := sm.Called(u)
	return args.Error(0).(*domain.Error)
}

func (sm *StorageMock) GetUserByEmail(email string) (*users.User, *domain.Error) {
	args := sm.Called(email)
	return args.Get(0).(*users.User), args.Error(1).(*domain.Error)
}

func (sm *StorageMock) GetUserIDByToken(token string) (*string, *domain.Error) {
	args := sm.Called(token)
	return args.Get(0).(*string), args.Error(1).(*domain.Error)
}
