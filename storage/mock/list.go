package mock

import "github.com/Sehsyha/crounch-back/model"

func (sm *StorageMock) CreateList(list *model.List) error {
	args := sm.Called(list)
	return args.Error(0)
}

func (sm *StorageMock) GetOwnerLists(ownerID string) ([]*model.List, error) {
	args := sm.Called(ownerID)
	return args.Get(0).([]*model.List), args.Error(1)
}
