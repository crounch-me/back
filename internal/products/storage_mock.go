package products

import (
	"github.com/crounch-me/back/internal"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (ps *StorageMock) CreateProduct(id, name, ownerID string) *internal.Error {
	args := ps.Called(id, name, ownerID)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*internal.Error)
}

func (ps *StorageMock) GetProduct(id string) (*Product, *internal.Error) {
	args := ps.Called(id)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*Product), nil
	}
	return nil, err.(*internal.Error)
}

func (ps *StorageMock) SearchDefaults(lowerCasedName string, id string) ([]*Product, *internal.Error) {
	args := ps.Called(lowerCasedName, id)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).([]*Product), nil
	}

	return nil, err.(*internal.Error)
}
