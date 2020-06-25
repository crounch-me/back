package products

import (
	"github.com/crounch-me/back/domain"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (ps *StorageMock) CreateProduct(id, name, ownerID string) *domain.Error {
	args := ps.Called(id, name, ownerID)
	err := args.Error(0)
	if err == nil {
		return nil
	}
	return err.(*domain.Error)
}

func (ps *StorageMock) GetProduct(id string) (*Product, *domain.Error) {
	args := ps.Called(id)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*Product), nil
	} else {
		return nil, err.(*domain.Error)
	}
}
