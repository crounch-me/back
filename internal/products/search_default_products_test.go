package products

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
	"github.com/stretchr/testify/assert"
)

func TestSearchDefaultProductsError(t *testing.T) {
	name := "name"

	productStorageMock := &StorageMock{}
	productStorageMock.On("SearchDefaults", name, users.AdminID).Return(nil, domain.NewError(domain.UnknownErrorCode))

	productService := &ProductService{
		ProductStorage: productStorageMock,
	}

	result, err := productService.SearchDefaults(name)

	assert.Empty(t, result)
	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestSearchDefaultProductsOK(t *testing.T) {
	name := "name"

	expectedProducts := []*Product{
		{
			ID: "id 1",
		},
		{
			ID: "id 2",
		},
	}

	productStorageMock := &StorageMock{}
	productStorageMock.On("SearchDefaults", name, users.AdminID).Return(expectedProducts, nil)

	productService := &ProductService{
		ProductStorage: productStorageMock,
	}

	result, err := productService.SearchDefaults(name)

	assert.Equal(t, expectedProducts, result)
	assert.Empty(t, err)
}
