package list

import (
	"time"

	"github.com/crounch-me/back/internal"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (sm *StorageMock) CreateList(id, name string, creationDate time.Time) *internal.Error {
	args := sm.Called(id, name, creationDate)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*internal.Error)
}

func (sm *StorageMock) ArchiveList(listID string, archivationDate time.Time) *internal.Error {
	args := sm.Called(listID, archivationDate)
	err := args.Error(0)
	if err != nil {
		return err.(*internal.Error)
	}

	return nil
}

func (sm *StorageMock) GetUsersLists(ownerID string) ([]*List, *internal.Error) {
	args := sm.Called(ownerID)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).([]*List), nil
	}

	return nil, err.(*internal.Error)
}

func (sm *StorageMock) GetList(id string) (*List, *internal.Error) {
	args := sm.Called(id)
	list := args.Get(0)
	if list == nil {
		return nil, args.Error(1).(*internal.Error)
	}

	return list.(*List), nil
}

func (sm *StorageMock) GetProductInList(productID string, listID string) (*ProductInListLink, *internal.Error) {
	args := sm.Called(productID, listID)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*ProductInListLink), nil
	}

	return nil, err.(*internal.Error)
}

func (sm *StorageMock) AddProductToList(productID string, listID string) *internal.Error {
	args := sm.Called(productID, listID)
	err := args.Error(0)
	if err == nil {
		return nil
	}

	return err.(*internal.Error)
}

func (sm *StorageMock) DeleteProductsFromList(listID string) *internal.Error {
	args := sm.Called(listID)
	err := args.Error(0)
	if err == nil {
		return nil
	}

	return err.(*internal.Error)
}

func (ps *StorageMock) DeleteProductFromList(productID string, listID string) *internal.Error {
	args := ps.Called(productID, listID)
	err := args.Error(0)
	if err == nil {
		return nil
	}

	return err.(*internal.Error)
}

func (sm *StorageMock) DeleteList(listID string) *internal.Error {
	args := sm.Called(listID)
	err := args.Error(0)
	if err == nil {
		return nil
	}

	return err.(*internal.Error)
}

func (sm *StorageMock) GetProductsOfList(listID string) ([]*ProductInList, *internal.Error) {
	args := sm.Called(listID)
	err := args.Error(1)
	if err != nil {
		return nil, err.(*internal.Error)
	}

	return args.Get(0).([]*ProductInList), nil
}

func (sm *StorageMock) UpdateProductInList(updateProductInList *UpdateProductInList, productID, listID string) (*ProductInListLink, *internal.Error) {
	args := sm.Called(updateProductInList, productID, listID)
	err := args.Error(1)

	if err != nil {
		return nil, err.(*internal.Error)
	}

	return args.Get(0).(*ProductInListLink), nil
}
