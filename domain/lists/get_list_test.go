package lists

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
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

func TestGetListUnauthorized(t *testing.T) {
	listID := "list-id"
	userID := "userID1"
	anotherUserID := "userID2"
	list := &List{
		Owner: &users.User{
			ID: anotherUserID,
		},
	}
	storageMock := &StorageMock{}

	storageMock.On("GetList", listID).Return(list, nil)

	listService := &ListService{
		ListStorage: storageMock,
	}
	result, err := listService.GetList(listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	assert.Empty(t, result)
	assert.Equal(t, domain.ForbiddenErrorCode, err.Code)
}

func TestGetListOK(t *testing.T) {
	listID := "list-id"
	userID := "user-id"
	list := &List{
		ID: listID,
		Owner: &users.User{
			ID: userID,
		},
	}
	storageMock := &StorageMock{}

	storageMock.On("GetList", listID).Return(list, nil)

	listService := &ListService{
		ListStorage: storageMock,
	}
	result, err := listService.GetList(listID, userID)

	storageMock.AssertCalled(t, "GetList", listID)
	assert.Equal(t, list, result)
	assert.Empty(t, err)
}
