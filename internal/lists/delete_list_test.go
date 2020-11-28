package lists

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/contributors"
	"github.com/stretchr/testify/assert"
)

func TestDeleteListGetListError(t *testing.T) {
	userID := "user-id"
	listID := "list-id"

  listStorageMock := &StorageMock{}

	listService := &ListService{
		ListStorage: listStorageMock,
  }

	listStorageMock.On("GetList", listID).Return(nil, domain.NewError(ListNotFoundErrorCode))

	err := listService.DeleteList(listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "DeleteProductsFromList")
	listStorageMock.AssertNotCalled(t, "DeleteList")
	assert.Equal(t, ListNotFoundErrorCode, err.Code)
}
func TestDeleteListGetContributorsIDsError(t *testing.T) {
	userID := "user-id"
	listID := "list-id"
	list := &List{
		ID: listID,
	}

  listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(list, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{}, domain.NewError(domain.UnknownErrorCode))

	listService := &ListService{
    ListStorage: listStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	err := listService.DeleteList(listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "DeleteProductsFromList")
  listStorageMock.AssertNotCalled(t, "DeleteList")
  contributorStorageMock.AssertCalled(t, "GetContributorsIDs", listID)

	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestDeleteListUnauthorized(t *testing.T) {
	userID := "user-id"
	anotherUserID := "another-user-id"
	listID := "list-id"
	list := &List{
		ID: listID,
	}

  listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(list, nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{anotherUserID}, nil)

	listService := &ListService{
    ListStorage: listStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	err := listService.DeleteList(listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertNotCalled(t, "DeleteProductsFromList")
  listStorageMock.AssertNotCalled(t, "DeleteList")
  contributorStorageMock.AssertCalled(t, "GetContributorsIDs", listID)

	assert.Equal(t, domain.ForbiddenErrorCode, err.Code)
}

func TestDeleteListDeleteProductsFromListError(t *testing.T) {
	userID := "user-id"
	listID := "list-id"
	list := &List{
		ID: listID,
	}

  listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(list, nil)
	listStorageMock.On("DeleteProductsFromList", listID).Return(domain.NewError(domain.UnknownErrorCode))

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	listService := &ListService{
    ListStorage: listStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	err := listService.DeleteList(listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertCalled(t, "DeleteProductsFromList", listID)
	listStorageMock.AssertNotCalled(t, "DeleteList")
	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestDeleteListDeleteListError(t *testing.T) {
	userID := "user-id"
	listID := "list-id"
	list := &List{
		ID: listID,
	}

  listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(list, nil)
	listStorageMock.On("DeleteProductsFromList", listID).Return(nil)
	listStorageMock.On("DeleteList", listID).Return(domain.NewError(domain.UnknownErrorCode))

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	listService := &ListService{
    ListStorage: listStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	err := listService.DeleteList(listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertCalled(t, "DeleteProductsFromList", listID)
	listStorageMock.AssertCalled(t, "DeleteList", listID)
	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestDeleteListOK(t *testing.T) {
	userID := "user-id"
	listID := "list-id"
	list := &List{
		ID: listID,
	}

	listStorageMock := &StorageMock{}
	listStorageMock.On("GetList", listID).Return(list, nil)
	listStorageMock.On("DeleteProductsFromList", listID).Return(nil)
  listStorageMock.On("DeleteList", listID).Return(nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("GetContributorsIDs", listID).Return([]string{userID}, nil)

	listService := &ListService{
    ListStorage: listStorageMock,
    ContributorStorage: contributorStorageMock,
	}

	err := listService.DeleteList(listID, userID)

	listStorageMock.AssertCalled(t, "GetList", listID)
	listStorageMock.AssertCalled(t, "DeleteProductsFromList", listID)
	listStorageMock.AssertCalled(t, "DeleteList", listID)
	assert.Empty(t, err)
}
