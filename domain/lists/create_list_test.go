package lists

import (
	"testing"
	"time"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/contributors"
	"github.com/crounch-me/back/domain/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestCreateListCreateListError(t *testing.T) {
	id := "id"
	name := "name"
	userID := "user-id"
	generationMock := &domain.GenerationMock{}
	generationMock.On("GenerateID").Return(id, nil)

	storageMock := &StorageMock{}
	storageMock.On("CreateList", id, name, mock.Anything).Return(domain.NewError(domain.UnknownErrorCode))

	listService := &ListService{
		Generation:  generationMock,
		ListStorage: storageMock,
	}

	result, err := listService.CreateList(name, userID)

	generationMock.AssertCalled(t, "GenerateID")
	storageMock.AssertCalled(t, "CreateList", id, name, mock.Anything)
	assert.Empty(t, result)
	assert.Equal(t, err.Code, domain.UnknownErrorCode)
}

func TestCreateListCreateContributorError(t *testing.T) {
	id := "id"
	name := "name"
	userID := "user-id"
	generationMock := &domain.GenerationMock{}
	generationMock.On("GenerateID").Return(id, nil)

	storageMock := &StorageMock{}
	storageMock.On("CreateList", id, name, mock.Anything).Return(nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("CreateContributor", id, userID).Return(domain.NewError(domain.UnknownErrorCode))

	listService := &ListService{
		Generation:  generationMock,
    ListStorage: storageMock,
    ContributorStorage: contributorStorageMock,
	}

	result, err := listService.CreateList(name, userID)

	generationMock.AssertCalled(t, "GenerateID")
	storageMock.AssertCalled(t, "CreateList", id, name, mock.Anything)

	assert.Empty(t, result)
	assert.Equal(t, err.Code, domain.UnknownErrorCode)
}

func TestCreateListOK(t *testing.T) {
	id := "id"
	name := "name"
	userID := "user-id"
	generationMock := &domain.GenerationMock{}
	generationMock.On("GenerateID").Return(id, nil)

	storageMock := &StorageMock{}
	storageMock.On("CreateList", id, name, mock.Anything).Return(nil)

  contributorStorageMock := &contributors.StorageMock{}
  contributorStorageMock.On("CreateContributor", id, userID).Return(nil)

	listService := &ListService{
		Generation:  generationMock,
    ListStorage: storageMock,
    ContributorStorage: contributorStorageMock,
	}

	result, err := listService.CreateList(name, userID)
	fakeCreationDate := time.Now()
	expectedList := &List{
		ID:   id,
		Name: name,
    CreationDate: fakeCreationDate,
    Contributors: []*users.User{
      {
        ID: userID,
      },
    },
	}

	generationMock.AssertCalled(t, "GenerateID")
	storageMock.AssertCalled(t, "CreateList", id, name, mock.Anything)

	result.CreationDate = fakeCreationDate

	assert.Equal(t, result, expectedList)
	assert.Empty(t, err)
}
