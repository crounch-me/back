package products

import (
	"testing"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/users"
	"github.com/stretchr/testify/assert"
)

func TestGetProductGetProductError(t *testing.T) {
	productID := "product-id"
	userID := "user-id"

	productStorageMock := &StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(nil, internal.NewError(ProductNotFoundErrorCode))

	productService := &ProductService{
		ProductStorage: productStorageMock,
	}

	result, err := productService.GetProduct(productID, userID)

	assert.Empty(t, result)
	assert.Equal(t, ProductNotFoundErrorCode, err.Code)
}

func TestGetProductUnauthorized(t *testing.T) {
	productID := "product-id"
	userID := "user-id"
	anotherUserID := "another-user-id"
	product := &Product{
		ID: productID,
		Owner: &users.User{
			ID: anotherUserID,
		},
	}

	productStorageMock := &StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(product, nil)

	productService := &ProductService{
		ProductStorage: productStorageMock,
	}

	result, err := productService.GetProduct(productID, userID)

	assert.Empty(t, result)
	assert.Equal(t, internal.ForbiddenErrorCode, err.Code)
}

func TestGetProductOK(t *testing.T) {
	productID := "product-id"
	userID := "user-id"
	product := &Product{
		ID: productID,
		Owner: &users.User{
			ID: userID,
		},
	}

	productStorageMock := &StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(product, nil)

	productService := &ProductService{
		ProductStorage: productStorageMock,
	}

	result, err := productService.GetProduct(productID, userID)

	assert.Equal(t, product, result)
	assert.Empty(t, err)
}
