package lists

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/contributors"
	"github.com/crounch-me/back/domain/products"
	"github.com/stretchr/testify/assert"
)

func TestDeleteProductFromListGetListError(t *testing.T) {
	userID := "user-id"
	productID := "product-id"
	listID := "list-id"

	productStorageMock := &products.StorageMock{}
  listStorageMock := &StorageMock{}
  listStorageMock.On("GetList", listID).Return(nil, domain.NewError(ListNotFoundErrorCode))

	listService := &ListService{
		ListStorage:    listStorageMock,
		ProductStorage: productStorageMock,
	}

	err := listService.DeleteProductFromList(productID, listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "DeleteProductFromList")

	productStorageMock.AssertNotCalled(t, "GetProduct")

	assert.Equal(t, ListNotFoundErrorCode, err.Code)
}

func TestDeleteProductFromListGetContributorsIDsError(t *testing.T) {
	userID := "user-id"
	productID := "product-id"
	listID := "list-id"
	list := &List{
		ID: listID,
	}

	productStorageMock := &products.StorageMock{}

  listStorageMock := &StorageMock{}
  listStorageMock.On("GetList", listID).Return(list, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{}, domain.NewError(domain.UnknownErrorCode))

	listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	err := listService.DeleteProductFromList(productID, listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "DeleteProductFromList")

	productStorageMock.AssertNotCalled(t, "GetProduct")

	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestDeleteProductFromListUnauthorized(t *testing.T) {
	userID := "user-id"
	productID := "product-id"
	anotherUserID := "another-user-id"
	listID := "list-id"
	list := &List{
		ID: listID,
	}

	productStorageMock := &products.StorageMock{}

  listStorageMock := &StorageMock{}
  listStorageMock.On("GetList", listID).Return(list, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{anotherUserID}, nil)

	listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

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
	}

	listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(list, nil)

  productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(nil, domain.NewError(products.ProductNotFoundErrorCode))

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

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
	}
	product := &products.Product{
		ID: productID,
	}

  listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(list, nil)
	listStorageMock.On("DeleteProductFromList", productID, listID).Return(domain.NewError(domain.UnknownErrorCode))

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(product, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

  listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	err := listService.DeleteProductFromList(productID, listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertCalled(t, "DeleteProductFromList", productID, listID)

	productStorageMock.AssertCalled(t, "GetProduct", productID)

	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}
