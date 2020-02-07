package mock

import "github.com/Sehsyha/crounch-back/model"

// CreateAuthorization mocks create authorization
func (sm *StorageMock) CreateAuthorization(u *model.User) (*model.Authorization, error) {
	args := sm.Called(u)
	return args.Get(0).(*model.Authorization), args.Error(1)
}
