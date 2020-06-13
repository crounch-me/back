package mock

import (
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/lists"
)

// CreateList mocks create list
func (sm *StorageMock) CreateList(list *lists.List) *domain.Error {
	args := sm.Called(list)
	return args.Error(0).(*domain.Error)
}

// GetOwnerLists mocks get list with owner id
func (sm *StorageMock) GetOwnerLists(ownerID string) ([]*lists.List, *domain.Error) {
	args := sm.Called(ownerID)
	return args.Get(0).([]*lists.List), args.Error(1).(*domain.Error)
}

// GetList mocks get list with an id
func (sm *StorageMock) GetList(id string) (*lists.List, *domain.Error) {
	args := sm.Called(id)
	return args.Get(0).(*lists.List), args.Error(1).(*domain.Error)
}

func (sm *StorageMock) GetProductInList(productID string, listID string) (*lists.ProductInList, *domain.Error) {
	args := sm.Called(productID, listID)
	return args.Get(0).(*lists.ProductInList), args.Error(1).(*domain.Error)
}

func (sm *StorageMock) AddProductToList(productID string, listID string) *domain.Error {
	args := sm.Called(productID, listID)
	return args.Error(0).(*domain.Error)
}
