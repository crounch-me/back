package lists

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/contributors"
	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddProductToListGetListError(t *testing.T) {
	listID := "list-id"
	userID := "user-id"
	listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(nil, domain.NewError(ListNotFoundErrorCode))

  contributorStorageMock := &contributors.StorageMock{}

	productStorageMock := &products.StorageMock{}

	listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	result, err := listService.AddProductToList("productID", listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "GetProductInList")
	listStorageMock.AssertNotCalled(t, "AddProductToList")
  contributorStorageMock.AssertNotCalled(t, "GetContributorsIDs")
	productStorageMock.AssertNotCalled(t, "GetProduct")

	assert.Empty(t, result)
	assert.Equal(t, err.Code, ListNotFoundErrorCode)
}

func TestAddProductToListGetContributorsIDsError(t *testing.T) {
	listID := "list-id"
	userID := "user-id"
	lists := &List{
		ID: listID,
	}
	listStorageMock := &StorageMock{}
  listStorageMock.On("GetList", listID).Return(lists, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{}, domain.NewError(domain.UnknownErrorCode))

	productStorageMock := &products.StorageMock{}

	listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	result, err := listService.AddProductToList("productID", listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "GetProductInList")
	listStorageMock.AssertNotCalled(t, "AddProductToList")
  contributorStorageMock.AssertCalled(t, "GetContributorsIDs", listID)
	productStorageMock.AssertNotCalled(t, "GetProduct")

	assert.Empty(t, result)
	assert.Equal(t, err.Code, domain.UnknownErrorCode)
}

func TestAddProductToListUserNotAuthorizedOnList(t *testing.T) {
	listID := "list-id"
	userID := "user-id"
	anotherUserID := "anotherUserID"
	lists := &List{
		ID: listID,
	}
	listStorageMock := &StorageMock{}
  listStorageMock.On("GetList", listID).Return(lists, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{anotherUserID}, nil)

	productStorageMock := &products.StorageMock{}

	listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	result, err := listService.AddProductToList("productID", listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "GetProductInList")
	listStorageMock.AssertNotCalled(t, "AddProductToList")
  contributorStorageMock.AssertCalled(t, "GetContributorsIDs", listID)
	productStorageMock.AssertNotCalled(t, "GetProduct")

	assert.Empty(t, result)
	assert.Equal(t, err.Code, domain.ForbiddenErrorCode)
}

func TestAddProductToListGetProductError(t *testing.T) {
	productID := "productID"
	listID := "listID"
	userID := "userID"
	list := &List{
		ID: listID,
	}
	listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", mock.Anything).Return(list, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", mock.Anything).Return(nil, domain.NewError(products.ProductNotFoundErrorCode))

	listService := &ListService{
		ListStorage:    listStorageMock,
		ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	result, err := listService.AddProductToList("productID", listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "GetProductInList")
  listStorageMock.AssertNotCalled(t, "AddProductToList")
  contributorStorageMock.AssertCalled(t, "GetContributorsIDs", listID)
	productStorageMock.AssertCalled(t, "GetProduct", productID)

	assert.Empty(t, result)
	assert.Equal(t, err.Code, products.ProductNotFoundErrorCode)
}

func TestAddProductToListUserNotAuthorizedOnProduct(t *testing.T) {
	productID := "productID"
	listID := "listID"
	userID := "userID"
	anotherUserID := "anotherUserID"
	list := &List{
		ID: listID,
	}
	product := &products.Product{
		ID: productID,
		Owner: &users.User{
			ID: anotherUserID,
		},
	}
	listStorageMock := &StorageMock{}
  listStorageMock.On("GetList", mock.Anything).Return(list, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", mock.Anything).Return(product, nil)

	listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	result, err := listService.AddProductToList("productID", listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "GetProductInList")
	listStorageMock.AssertNotCalled(t, "AddProductToList")
  contributorStorageMock.AssertCalled(t, "GetContributorsIDs", listID)
	productStorageMock.AssertCalled(t, "GetProduct", productID)

	assert.Empty(t, result)
	assert.Equal(t, domain.ForbiddenErrorCode, err.Code)
}

func TestAddProductToListGetProductInListError(t *testing.T) {
	productID := "productID"
	listID := "listID"
	userID := "userID"
	list := &List{
		ID: listID,
	}
	product := &products.Product{
		ID: productID,
		Owner: &users.User{
			ID: userID,
		},
	}
	listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(list, nil)
	listStorageMock.On("GetProductInList", productID, listID).Return(nil, domain.NewError(domain.UnknownErrorCode))

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(product, nil)

	listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	result, err := listService.AddProductToList("productID", listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertCalled(t, "GetProductInList", productID, listID)
	listStorageMock.AssertNotCalled(t, "AddProductToList")
  contributorStorageMock.AssertCalled(t, "GetContributorsIDs", listID)
	productStorageMock.AssertCalled(t, "GetProduct", productID)

	assert.Empty(t, result)
	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestAddProductToListDuplicatedProductInListError(t *testing.T) {
	productID := "productID"
	listID := "listID"
	userID := "userID"
	list := &List{
		ID: listID,
	}
	product := &products.Product{
		ID: productID,
		Owner: &users.User{
			ID: userID,
		},
	}
	productInList := &ProductInListLink{
		ListID:    listID,
		ProductID: productID,
	}
	listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(list, nil)
	listStorageMock.On("GetProductInList", productID, listID).Return(productInList, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(product, nil)

	listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	result, err := listService.AddProductToList("productID", listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertCalled(t, "GetProductInList", productID, listID)
  listStorageMock.AssertNotCalled(t, "AddProductToList")
  contributorStorageMock.AssertCalled(t, "GetContributorsIDs", listID)
	productStorageMock.AssertCalled(t, "GetProduct", productID)

	assert.Empty(t, result)
	assert.Equal(t, DuplicateProductInListErrorCode, err.Code)
}

func TestAddProductToListAddProductToListError(t *testing.T) {
	productID := "productID"
	listID := "listID"
	userID := "userID"
	list := &List{
		ID: listID,
	}
	product := &products.Product{
		ID: productID,
		Owner: &users.User{
			ID: userID,
		},
	}
	listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(list, nil)
	listStorageMock.On("GetProductInList", productID, listID).Return(nil, domain.NewError(ProductInListNotFoundErrorCode))
	listStorageMock.On("AddProductToList", productID, listID).Return(domain.NewError(domain.UnknownErrorCode))

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(product, nil)

	listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	result, err := listService.AddProductToList("productID", listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertCalled(t, "GetProductInList", productID, listID)
  listStorageMock.AssertCalled(t, "AddProductToList", productID, listID)
  contributorStorageMock.AssertCalled(t, "GetContributorsIDs", listID)
	productStorageMock.AssertCalled(t, "GetProduct", productID)

	assert.Empty(t, result)
	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestAddProductToListOK(t *testing.T) {
	productID := "productID"
	listID := "listID"
	userID := "userID"
	list := &List{
		ID: listID,
	}
	product := &products.Product{
		ID: productID,
		Owner: &users.User{
			ID: userID,
		},
	}
	listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(list, nil)
	listStorageMock.On("GetProductInList", productID, listID).Return(nil, domain.NewError(ProductInListNotFoundErrorCode))
	listStorageMock.On("AddProductToList", productID, listID).Return(nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	productStorageMock := &products.StorageMock{}
	productStorageMock.On("GetProduct", productID).Return(product, nil)

	listService := &ListService{
		ListStorage:    listStorageMock,
    ProductStorage: productStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	result, err := listService.AddProductToList("productID", listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertCalled(t, "GetProductInList", productID, listID)
  listStorageMock.AssertCalled(t, "AddProductToList", productID, listID)
  contributorStorageMock.AssertCalled(t, "GetContributorsIDs", listID)
	productStorageMock.AssertCalled(t, "GetProduct", productID)

	expectedProductInList := &ProductInListLink{
		ListID:    listID,
		ProductID: productID,
	}

	assert.Equal(t, expectedProductInList, result)
	assert.Empty(t, err)
}
