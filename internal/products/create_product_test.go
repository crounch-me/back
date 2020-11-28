package products

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductGenerateIDError(t *testing.T) {
	generationMock := &domain.GenerationMock{}
	generationMock.On("GenerateID").Return("", domain.NewError(domain.UnknownErrorCode))

	productService := &ProductService{
		Generation: generationMock,
	}

	result, err := productService.CreateProduct("name", "user-id")

	assert.Empty(t, result)
	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestCreateProductCreateProductError(t *testing.T) {
	name := "name"
	userID := "user-id"
	productID := "product-id"

	generationMock := &domain.GenerationMock{}
	generationMock.On("GenerateID").Return(productID, nil)

	productStorageMock := &StorageMock{}
	productStorageMock.On("CreateProduct", productID, name, userID).Return(domain.NewError(domain.UnknownErrorCode))

	productService := &ProductService{
		Generation:     generationMock,
		ProductStorage: productStorageMock,
	}

	result, err := productService.CreateProduct(name, userID)

	assert.Empty(t, result)
	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestCreateProductOK(t *testing.T) {
	name := "name"
	userID := "user-id"
	productID := "product-id"

	generationMock := &domain.GenerationMock{}
	generationMock.On("GenerateID").Return(productID, nil)

	productStorageMock := &StorageMock{}
	productStorageMock.On("CreateProduct", productID, name, userID).Return(nil)

	productService := &ProductService{
		Generation:     generationMock,
		ProductStorage: productStorageMock,
	}

	expectedProduct := &Product{
		ID:   productID,
		Name: name,
		Owner: &users.User{
			ID: userID,
		},
	}

	result, err := productService.CreateProduct(name, userID)

	assert.Equal(t, expectedProduct, result)
	assert.Empty(t, err)
}
