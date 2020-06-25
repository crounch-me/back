package lists

import (
	"time"

	"github.com/crounch-me/back/domain"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (sm *StorageMock) CreateList(id, name, ownerID string, creationDate time.Time) *domain.Error {
	args := sm.Called(id, name, ownerID, creationDate)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*domain.Error)
}

func (sm *StorageMock) GetOwnersLists(ownerID string) ([]*List, *domain.Error) {
	args := sm.Called(ownerID)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).([]*List), nil
	}
	return nil, err.(*domain.Error)
}

func (sm *StorageMock) GetList(id string) (*List, *domain.Error) {
	args := sm.Called(id)
	list := args.Get(0)
	if list == nil {
		return nil, args.Error(1).(*domain.Error)
	} else {
		return list.(*List), nil
	}
}

func (sm *StorageMock) GetProductInList(productID string, listID string) (*ProductInList, *domain.Error) {
	args := sm.Called(productID, listID)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*ProductInList), nil
	}
	return nil, err.(*domain.Error)
}

func (sm *StorageMock) AddProductToList(productID string, listID string) *domain.Error {
	args := sm.Called(productID, listID)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*domain.Error)
}

func (sm *StorageMock) DeleteProductsFromList(listID string) *domain.Error {
	args := sm.Called(listID)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*domain.Error)
}

func (sm *StorageMock) DeleteList(listID string) *domain.Error {
	args := sm.Called(listID)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*domain.Error)
}
