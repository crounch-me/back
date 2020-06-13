package mock

import (
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/products"
)

// CreateProduct mocks create product
func (sm *StorageMock) CreateProduct(product *products.Product) *domain.Error {
	args := sm.Called(product)
	return args.Error(0).(*domain.Error)
}

// GetProduct mocks get product with an id
func (sm *StorageMock) GetProduct(id string) (*products.Product, *domain.Error) {
	args := sm.Called(id)
	return args.Get(0).(*products.Product), args.Error(1).(*domain.Error)
}
