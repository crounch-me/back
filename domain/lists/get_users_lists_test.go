package lists

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUsersListsStorageError(t *testing.T) {
	storageMock := &StorageMock{}
	storageMock.On("GetUsersLists", mock.Anything).Return(nil, domain.NewError(domain.UnknownErrorCode))

	listService := &ListService{
		ListStorage: storageMock,
	}

	result, err := listService.GetUsersLists("userID")

	storageMock.AssertCalled(t, "GetUsersLists", mock.Anything)
	assert.Empty(t, result)
	assert.Equal(t, err.Code, domain.UnknownErrorCode)
}

func TestGetUsersListStorageOK(t *testing.T) {
	lists := []*List{
		{},
	}
	storageMock := &StorageMock{}
	storageMock.On("GetUsersLists", mock.Anything).Return(lists, nil)

	listService := &ListService{
		ListStorage: storageMock,
	}

	result, err := listService.GetUsersLists("userID")

	storageMock.AssertCalled(t, "GetUsersLists", mock.Anything)
	assert.Equal(t, result, lists)
	assert.Empty(t, err)
}
