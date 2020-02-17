package mock

import "github.com/Sehsyha/crounch-back/model"

// CreateProduct mocks create product
func (sm *StorageMock) CreateProduct(product *model.Product) error {
	args := sm.Called(product)
	return args.Error(0)
}
