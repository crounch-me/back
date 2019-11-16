package mock

import "github.com/Sehsyha/crounch-back/model"

func (sm *StorageMock) CreateList(list *model.List) error {
	args := sm.Called(list)
	return args.Error(0)
}
