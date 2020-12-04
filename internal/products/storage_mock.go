package products

import (
	"github.com/crounch-me/back/internal/common/errors"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (ps *StorageMock) CreateProduct(id, name, ownerID string) *errors.Error {
	args := ps.Called(id, name, ownerID)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*errors.Error)
}

func (ps *StorageMock) GetProduct(id string) (*Product, *errors.Error) {
	args := ps.Called(id)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*Product), nil
	}
	return nil, err.(*errors.Error)
}

func (ps *StorageMock) SearchDefaults(lowerCasedName string, id string) ([]*Product, *errors.Error) {
	args := ps.Called(lowerCasedName, id)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).([]*Product), nil
	}

	return nil, err.(*errors.Error)
}
