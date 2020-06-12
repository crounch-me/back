package mock

import "github.com/crounch-me/back/model"

// CreateProduct mocks create product
func (sm *StorageMock) CreateProduct(product *model.Product) error {
	args := sm.Called(product)
	return args.Error(0)
}

// GetProduct mocks get product with an id
func (sm *StorageMock) GetProduct(id string) (*model.Product, error) {
	args := sm.Called(id)
	return args.Get(0).(*model.Product), args.Error(1)
}
