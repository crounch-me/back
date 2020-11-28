package products

import (
	"testing"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/users"
	"github.com/stretchr/testify/assert"
)

func TestSearchDefaultProductsError(t *testing.T) {
	name := "name"

	productStorageMock := &StorageMock{}
	productStorageMock.On("SearchDefaults", name, users.AdminID).Return(nil, internal.NewError(internal.UnknownErrorCode))

	productService := &ProductService{
		ProductStorage: productStorageMock,
	}

	result, err := productService.SearchDefaults(name)

	assert.Empty(t, result)
	assert.Equal(t, internal.UnknownErrorCode, err.Code)
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
