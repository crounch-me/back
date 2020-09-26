package lists

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
	"github.com/stretchr/testify/assert"
)

func TestDeleteProductFromListGetListError(t *testing.T) {
	userID := "user-id"
	productID := "product-id"
	listID := "list-id"

	productStorageMock := &products.StorageMock{}
	listStorageMock := &StorageMock{}
	listService := &ListService{
		ListStorage:    listStorageMock,
		ProductStorage: productStorageMock,
	}

	listStorageMock.On("GetList", listID).Return(nil, domain.NewError(ListNotFoundErrorCode))

	err := listService.DeleteProductFromList(productID, listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "DeleteProductFromList")

	productStorageMock.AssertNotCalled(t, "GetProduct")

	assert.Equal(t, ListNotFoundErrorCode, err.Code)
}

func TestDeleteProductFromListUnauthorized(t *testing.T) {
	userID := "user-id"
	productID := "product-id"
	anotherUserID := "another-user-id"
	listID := "list-id"
	list := &List{
		ID: listID,
		Owner: &users.User{
			ID: anotherUserID,
		},
	}

	productStorageMock := &products.StorageMock{}
	listStorageMock := &StorageMock{}
	listService := &ListService{
		ListStorage:    listStorageMock,
		ProductStorage: productStorageMock,
	}

	listStorageMock.On("GetList", listID).Return(list, nil)

	err := listService.DeleteProductFromList(productID, listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "DeleteProductFromList")

	productStorageMock.AssertNotCalled(t, "GetProduct")

	assert.Equal(t, domain.ForbiddenErrorCode, err.Code)
}

func TestDeleteProductFromListGetProductError(t *testing.T) {
	userID := "user-id"
	productID := "product-id"
	listID := "list-id"
	list := &List{
		ID: listID,
		Owner: &users.User{
			ID: userID,
		},
	}

	listStorageMock := &StorageMock{}
	productStorageMock := &products.StorageMock{}
	listService := &ListService{
		ListStorage:    listStorageMock,
		ProductStorage: productStorageMock,
	}

	listStorageMock.On("GetList", listID).Return(list, nil)
	productStorageMock.On("GetProduct", productID).Return(nil, domain.NewError(products.ProductNotFoundErrorCode))

	err := listService.DeleteProductFromList(productID, listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "DeleteProductFromList")

	productStorageMock.AssertCalled(t, "GetProduct", productID)

	assert.Equal(t, products.ProductNotFoundErrorCode, err.Code)
}

func TestDeleteProductFromListDeleteProductFromListError(t *testing.T) {
	userID := "user-id"
	productID := "product-id"
	listID := "list-id"
	list := &List{
		ID: listID,
		Owner: &users.User{
			ID: userID,
		},
	}
	product := &products.Product{
		ID: productID,
	}

	listStorageMock := &StorageMock{}
	productStorageMock := &products.StorageMock{}
	listService := &ListService{
		ListStorage:    listStorageMock,
		ProductStorage: productStorageMock,
	}

	listStorageMock.On("GetList", listID).Return(list, nil)
	listStorageMock.On("DeleteProductFromList", productID, listID).Return(domain.NewError(domain.UnknownErrorCode))

	productStorageMock.On("GetProduct", productID).Return(product, nil)

	err := listService.DeleteProductFromList(productID, listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertCalled(t, "DeleteProductFromList", productID, listID)

	productStorageMock.AssertCalled(t, "GetProduct", productID)

	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}
