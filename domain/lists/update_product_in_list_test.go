package lists

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/contributors"
	"github.com/crounch-me/back/domain/products"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProductInListGetListError(t *testing.T) {
	listID := "listID"
	productID := "productID"
	userID := "userID"

	updateProductInList := &UpdateProductInList{
		Buyed: true,
	}

	storageMock := &StorageMock{}
	storageMock.On("GetList", listID).Return(nil, domain.NewError(ListNotFoundErrorCode))

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(nil, domain.NewError(products.ProductNotFoundErrorCode))

	listService := &ListService{
		ListStorage:    storageMock,
		ProductStorage: productStorageMock,
	}

	result, err := listService.UpdateProductInList(updateProductInList, productID, listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	storageMock.AssertNotCalled(t, "UpdateProductInList")
	productStorageMock.AssertNotCalled(t, "GetProduct")
	assert.Empty(t, result)
	assert.Equal(t, err.Code, ListNotFoundErrorCode)
}

func TestUpdateProductInListGetContributorsIDsError(t *testing.T) {
	listID := "listID"
	productID := "productID"
	userID := "userID"

	updateProductInList := &UpdateProductInList{
		Buyed: true,
	}

	list := &List{}

	storageMock := &StorageMock{}
	storageMock.On("GetList", listID).Return(list, nil)

	contributorStorageMock := &contributors.StorageMock{}
	contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{}, domain.NewError(domain.UnknownErrorCode))

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(nil, domain.NewError(products.ProductNotFoundErrorCode))

	listService := &ListService{
		ListStorage:        storageMock,
		ProductStorage:     productStorageMock,
		ContributorStorage: contributorStorageMock,
	}

	result, err := listService.UpdateProductInList(updateProductInList, productID, listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	storageMock.AssertNotCalled(t, "UpdateProductInList")
	productStorageMock.AssertNotCalled(t, "GetProduct")
	assert.Empty(t, result)
	assert.Equal(t, err.Code, domain.UnknownErrorCode)
}

func TestUpdateProductInListUserNotAuthorized(t *testing.T) {
	listID := "listID"
	productID := "productID"
	userID := "userID"
	anotherUserID := "anotherUserID"

	updateProductInList := &UpdateProductInList{
		Buyed: true,
	}

	list := &List{}

	storageMock := &StorageMock{}
	storageMock.On("GetList", listID).Return(list, nil)

	contributorStorageMock := &contributors.StorageMock{}
	contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{anotherUserID}, nil)

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(nil, domain.NewError(products.ProductNotFoundErrorCode))

	listService := &ListService{
		ListStorage:        storageMock,
		ProductStorage:     productStorageMock,
		ContributorStorage: contributorStorageMock,
	}

	result, err := listService.UpdateProductInList(updateProductInList, productID, listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	storageMock.AssertNotCalled(t, "UpdateProductInList")
	productStorageMock.AssertNotCalled(t, "GetProduct")
	assert.Empty(t, result)
	assert.Equal(t, err.Code, domain.ForbiddenErrorCode)
}

func TestUpdateProductInListGetProductError(t *testing.T) {
	listID := "listID"
	productID := "productID"
	userID := "userID"

	updateProductInList := &UpdateProductInList{
		Buyed: true,
	}

	list := &List{}

	storageMock := &StorageMock{}
	storageMock.On("GetList", listID).Return(list, nil)

	contributorStorageMock := &contributors.StorageMock{}
	contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(nil, domain.NewError(products.ProductNotFoundErrorCode))

	listService := &ListService{
		ListStorage:        storageMock,
		ProductStorage:     productStorageMock,
		ContributorStorage: contributorStorageMock,
	}

	result, err := listService.UpdateProductInList(updateProductInList, productID, listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	storageMock.AssertNotCalled(t, "UpdateProductInList")
	productStorageMock.AssertCalled(t, "GetProduct", productID)
	assert.Empty(t, result)
	assert.Equal(t, err.Code, products.ProductNotFoundErrorCode)
}

func TestUpdateProductInListUpdateProductInListError(t *testing.T) {
	listID := "listID"
	productID := "productID"
	userID := "userID"

	updateProductInList := &UpdateProductInList{
		Buyed: true,
	}

	list := &List{}

	product := &products.Product{}

	storageMock := &StorageMock{}
	storageMock.On("GetList", listID).Return(list, nil)
	storageMock.On("UpdateProductInList", updateProductInList, productID, listID).Return(nil, domain.NewError(ProductInListNotFoundErrorCode))

	contributorStorageMock := &contributors.StorageMock{}
	contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(product, nil)

	listService := &ListService{
		ListStorage:        storageMock,
		ProductStorage:     productStorageMock,
		ContributorStorage: contributorStorageMock,
	}

	result, err := listService.UpdateProductInList(updateProductInList, productID, listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	storageMock.AssertCalled(t, "UpdateProductInList", updateProductInList, productID, listID)
	productStorageMock.AssertCalled(t, "GetProduct", productID)
	assert.Empty(t, result)
	assert.Equal(t, err.Code, ProductInListNotFoundErrorCode)
}

func TestUpdateProductInListOK(t *testing.T) {
	listID := "listID"
	productID := "productID"
	userID := "userID"
	buyed := true

	updateProductInList := &UpdateProductInList{
		Buyed: buyed,
	}

	list := &List{}

	product := &products.Product{}

	productInList := &ProductInListLink{
		ListID:    listID,
		ProductID: productID,
		Buyed:     buyed,
	}

	storageMock := &StorageMock{}
	storageMock.On("GetList", listID).Return(list, nil)
	storageMock.On("UpdateProductInList", updateProductInList, productID, listID).Return(productInList, nil)

	contributorStorageMock := &contributors.StorageMock{}
	contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(product, nil)

	listService := &ListService{
		ListStorage:    storageMock,
		ProductStorage: productStorageMock,
		ContributorStorage: contributorStorageMock,
	}

	result, err := listService.UpdateProductInList(updateProductInList, productID, listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	storageMock.AssertCalled(t, "UpdateProductInList", updateProductInList, productID, listID)
	productStorageMock.AssertCalled(t, "GetProduct", productID)
	assert.Equal(t, result, productInList)
	assert.Empty(t, err)
}
