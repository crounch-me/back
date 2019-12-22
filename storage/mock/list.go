package mock

import "github.com/Sehsyha/crounch-back/model"

// CreateList mocks create list
func (sm *StorageMock) CreateList(list *model.List) error {
	args := sm.Called(list)
	return args.Error(0)
}

// GetOwnerLists mocks get list with owner id
func (sm *StorageMock) GetOwnerLists(ownerID string) ([]*model.List, error) {
	args := sm.Called(ownerID)
	return args.Get(0).([]*model.List), args.Error(1)
}

// GetList mocks get list with an id
func (sm *StorageMock) GetList(id string) (*model.List, error) {
	args := sm.Called(id)
	return args.Get(0).(*model.List), args.Error(1)
}

// AddOFFProductToList mocks add a product from open food facts to a list
func (sm *StorageMock) AddOFFProductToList(listID string, offProduct *model.OFFProduct) error {
	args := sm.Called(listID, offProduct)
	return args.Error(0)
}

// GetProdutsFromList mocks the retrieval of OFFProducts from a list
func (sm *StorageMock) GetOFFProducts(listID string) ([]*model.OFFProduct, error) {
	args := sm.Called(listID)
	return args.Get(0).([]*model.OFFProduct), args.Error(1)
}
