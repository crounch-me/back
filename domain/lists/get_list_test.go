package lists

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/contributors"
	"github.com/crounch-me/back/domain/products"
	"github.com/stretchr/testify/assert"
)

func TestGetListStorageError(t *testing.T) {
	listID := "list-id"
	storageMock := &StorageMock{}
	storageMock.On("GetList", listID).Return(nil, domain.NewError(ListNotFoundErrorCode))

	listService := &ListService{
		ListStorage: storageMock,
	}

	result, err := listService.GetList(listID, "userID")

	storageMock.AssertCalled(t, "GetList", listID)
	assert.Empty(t, result)
	assert.Equal(t, err.Code, ListNotFoundErrorCode)
}

func TestGetListGetContributorsIDsError(t *testing.T) {
	listID := "list-id"
	userID := "userID1"
  list := &List{}

	storageMock := &StorageMock{}
	storageMock.On("GetList", listID).Return(list, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{}, domain.NewError(domain.UnknownErrorCode))

	listService := &ListService{
    ListStorage: storageMock,
    ContributorStorage: contributorStorageMock,
  }

	result, err := listService.GetList(listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	assert.Empty(t, result)
	assert.Equal(t, domain.UnknownErrorCode, err.Code)}

func TestGetListUnauthorized(t *testing.T) {
	listID := "list-id"
	userID := "userID1"
	anotherUserID := "userID2"
  list := &List{}

	storageMock := &StorageMock{}
	storageMock.On("GetList", listID).Return(list, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{anotherUserID}, nil)

	listService := &ListService{
    ListStorage: storageMock,
    ContributorStorage: contributorStorageMock,
  }

	result, err := listService.GetList(listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	assert.Empty(t, result)
	assert.Equal(t, domain.ForbiddenErrorCode, err.Code)
}

func TestGetListGetProductsOfListError(t *testing.T) {
	listID := "list-id"
	userID := "user-id"
	list := &List{
		ID: listID,
	}
	storageMock := &StorageMock{}

	storageMock.On("GetList", listID).Return(list, nil)
  storageMock.On("GetProductsOfList", listID).Return(nil, domain.NewError(domain.UnknownErrorCode))

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	listService := &ListService{
    ListStorage: storageMock,
    ContributorStorage: contributorStorageMock,
	}
	result, err := listService.GetList(listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	assert.Empty(t, result)
	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestGetListOK(t *testing.T) {
	listID := "list-id"
	userID := "user-id"
	productID := "product-id"
	productName := "product-name"
	listProducts := []*ProductInList{
		{
			Product: &products.Product{
				ID:   productID,
				Name: productName,
			},
		},
	}
	list := &List{
		ID: listID,
		Products: listProducts,
	}
	storageMock := &StorageMock{}

	storageMock.On("GetList", listID).Return(list, nil)
	storageMock.On("GetProductsOfList", listID).Return(listProducts, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	listService := &ListService{
    ListStorage: storageMock,
    ContributorStorage: contributorStorageMock,
	}
	result, err := listService.GetList(listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	assert.Equal(t, list, result)
	assert.Empty(t, err)
}
