package mock

import (
	"github.com/Sehsyha/crounch-back/model"
)

func (sm *StorageMock) CreateUser(u *model.User) error {
	args := sm.Called(u)
	return args.Error(0)
}

func (sm *StorageMock) GetUserByEmail(email string) (*model.User, error) {
	args := sm.Called(email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (sm *StorageMock) GetUserIDByToken(token string) (*string, error) {
	args := sm.Called(token)
	return args.Get(0).(*string), args.Error(1)
}
