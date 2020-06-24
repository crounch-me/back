package lists

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetOwnersListStorageError(t *testing.T) {
	storageMock := &StorageMock{}
	storageMock.On("GetOwnersLists", mock.Anything).Return(nil, domain.NewError(ListNotFoundErrorCode))

	listService := &ListService{
		ListStorage: storageMock,
	}

	result, err := listService.GetOwnersLists("userID")

	storageMock.AssertCalled(t, "GetOwnersLists", mock.Anything)
	assert.Empty(t, result)
	assert.Equal(t, err.Code, ListNotFoundErrorCode)
}

func TestGetOwnersListStorageOK(t *testing.T) {
	lists := []*List{
		{},
	}
	storageMock := &StorageMock{}
	storageMock.On("GetOwnersLists", mock.Anything).Return(lists, nil)

	listService := &ListService{
		ListStorage: storageMock,
	}

	result, err := listService.GetOwnersLists("userID")

	storageMock.AssertCalled(t, "GetOwnersLists", mock.Anything)
	assert.Equal(t, result, lists)
	assert.Empty(t, err)
}
