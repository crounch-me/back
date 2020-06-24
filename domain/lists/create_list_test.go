package lists

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
	"github.com/stretchr/testify/assert"
)

func TestGenerateIDError(t *testing.T) {
	generationMock := &domain.GenerationMock{}
	generationMock.On("GenerateID").Return("", domain.NewError(domain.UnknownErrorCode))

	storageMock := &StorageMock{}

	listService := &ListService{
		Generation:  generationMock,
		ListStorage: storageMock,
	}

	result, err := listService.CreateList("list", "userID")

	generationMock.AssertCalled(t, "GenerateID")
	storageMock.AssertNotCalled(t, "CreateList")
	assert.Empty(t, result)
	assert.Equal(t, err.Code, domain.UnknownErrorCode)
}

func TestServiceCreateListError(t *testing.T) {
	id := "id"
	name := "name"
	userID := "user-id"
	generationMock := &domain.GenerationMock{}
	generationMock.On("GenerateID").Return(id, nil)

	storageMock := &StorageMock{}
	storageMock.On("CreateList", id, name, userID).Return(domain.NewError(domain.UnknownErrorCode))

	listService := &ListService{
		Generation:  generationMock,
		ListStorage: storageMock,
	}

	result, err := listService.CreateList(name, userID)

	generationMock.AssertCalled(t, "GenerateID")
	storageMock.AssertCalled(t, "CreateList", id, name, userID)
	assert.Empty(t, result)
	assert.Equal(t, err.Code, domain.UnknownErrorCode)
}

func TestServiceCreateListOK(t *testing.T) {
	id := "id"
	name := "name"
	userID := "user-id"
	generationMock := &domain.GenerationMock{}
	generationMock.On("GenerateID").Return(id, nil)

	storageMock := &StorageMock{}
	storageMock.On("CreateList", id, name, userID).Return(nil)

	listService := &ListService{
		Generation:  generationMock,
		ListStorage: storageMock,
	}

	result, err := listService.CreateList(name, userID)
	expectedList := &List{
		ID:   id,
		Name: name,
		Owner: &users.User{
			ID: userID,
		},
	}

	generationMock.AssertCalled(t, "GenerateID")
	storageMock.AssertCalled(t, "CreateList", id, name, userID)
	assert.Equal(t, result, expectedList)
	assert.Empty(t, err)
}
