package lists

import (
	"time"

	"github.com/crounch-me/back/domain"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (sm *StorageMock) CreateList(id, name string, creationDate time.Time) *domain.Error {
	args := sm.Called(id, name, creationDate)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*domain.Error)
}

func (sm *StorageMock) ArchiveList(listID string, archivationDate time.Time) *domain.Error {
	args := sm.Called(listID, archivationDate)
	err := args.Error(0)
	if err != nil {
		return err.(*domain.Error)
	}

	return nil
}

func (sm *StorageMock) GetUsersLists(ownerID string) ([]*List, *domain.Error) {
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
	}

	return list.(*List), nil
}

func (sm *StorageMock) GetProductInList(productID string, listID string) (*ProductInListLink, *domain.Error) {
	args := sm.Called(productID, listID)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*ProductInListLink), nil
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

func (ps *StorageMock) DeleteProductFromList(productID string, listID string) *domain.Error {
	args := ps.Called(productID, listID)
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

func (sm *StorageMock) GetProductsOfList(listID string) ([]*ProductInList, *domain.Error) {
	args := sm.Called(listID)
	err := args.Error(1)
	if err != nil {
		return nil, err.(*domain.Error)
	}

	return args.Get(0).([]*ProductInList), nil
}

func (sm *StorageMock) UpdateProductInList(updateProductInList *UpdateProductInList, productID, listID string) (*ProductInListLink, *domain.Error) {
	args := sm.Called(updateProductInList, productID, listID)
	err := args.Error(1)

	if err != nil {
		return nil, err.(*domain.Error)
	}

	return args.Get(0).(*ProductInListLink), nil
}
