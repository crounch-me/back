package authorization

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuthorizationGetUserByEmailError(t *testing.T) {
	email := "email"
	password := "password"

	userStorageMock := &users.StorageMock{}
	userStorageMock.On("GetByEmail", email).Return(nil, domain.NewError(users.UserNotFoundErrorCode))

	authorizationService := &AuthorizationService{
		UserStorage: userStorageMock,
	}

	result, err := authorizationService.CreateAuthorization(email, password)

	userStorageMock.AssertCalled(t, "GetByEmail", email)

	assert.Empty(t, result)
	assert.Equal(t, users.UserNotFoundErrorCode, err.Code)
}

func TestCreateAuthorizationWrongPassword(t *testing.T) {
	email := "email"
	password := "wrong-password"
	userPassword := "password"
	user := &users.User{
		Password: &userPassword,
	}

	generationMock := &domain.GenerationMock{}
	generationMock.On("ComparePassword", userPassword, password).Return(false)

	userStorageMock := &users.StorageMock{}
	userStorageMock.On("GetByEmail", email).Return(user, nil)

	authorizationService := &AuthorizationService{
		UserStorage: userStorageMock,
		Generation:  generationMock,
	}

	result, err := authorizationService.CreateAuthorization(email, password)

	generationMock.AssertCalled(t, "ComparePassword", userPassword, password)
	generationMock.AssertNotCalled(t, "GenerateToken")
	userStorageMock.AssertCalled(t, "GetByEmail", email)
	assert.Empty(t, result)
	assert.Equal(t, WrongPasswordErrorCode, err.Code)
}

func TestCreateAuthorizationGenerateTokenError(t *testing.T) {
	email := "email"
	password := "wrong-password"
	userID := "user-id"
	token := "generated-token"
	user := &users.User{
		ID:       userID,
		Password: &password,
	}

	generationMock := &domain.GenerationMock{}
	generationMock.On("ComparePassword", password, password).Return(true)
	generationMock.On("GenerateToken").Return("", domain.NewError(domain.UnknownErrorCode))

	userStorageMock := &users.StorageMock{}
	userStorageMock.On("GetByEmail", email).Return(user, nil)

	authorizationStorageMock := &StorageMock{}
	authorizationStorageMock.On("CreateAuthorization", userID, token).Return(domain.NewError(domain.UnknownErrorCode))

	authorizationService := &AuthorizationService{
		AuthorizationStorage: authorizationStorageMock,
		UserStorage:          userStorageMock,
		Generation:           generationMock,
	}

	result, err := authorizationService.CreateAuthorization(email, password)

	generationMock.AssertCalled(t, "ComparePassword", password, password)
	generationMock.AssertCalled(t, "GenerateToken")
	userStorageMock.AssertCalled(t, "GetByEmail", email)
	authorizationStorageMock.AssertNotCalled(t, "CreateAuthorization")
	assert.Empty(t, result)
	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}
func TestCreateAuthorizationCreateAuthorizationError(t *testing.T) {
	email := "email"
	password := "wrong-password"
	userID := "user-id"
	token := "generated-token"
	user := &users.User{
		ID:       userID,
		Password: &password,
	}

	generationMock := &domain.GenerationMock{}
	generationMock.On("ComparePassword", password, password).Return(true)
	generationMock.On("GenerateToken").Return(token, nil)

	userStorageMock := &users.StorageMock{}
	userStorageMock.On("GetByEmail", email).Return(user, nil)

	authorizationStorageMock := &StorageMock{}
	authorizationStorageMock.On("CreateAuthorization", userID, token).Return(domain.NewError(domain.UnknownErrorCode))

	authorizationService := &AuthorizationService{
		AuthorizationStorage: authorizationStorageMock,
		UserStorage:          userStorageMock,
		Generation:           generationMock,
	}

	result, err := authorizationService.CreateAuthorization(email, password)

	generationMock.AssertCalled(t, "ComparePassword", password, password)
	generationMock.AssertCalled(t, "GenerateToken")
	userStorageMock.AssertCalled(t, "GetByEmail", email)
	authorizationStorageMock.AssertCalled(t, "CreateAuthorization", userID, token)
	assert.Empty(t, result)
	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestCreateAuthorizationOK(t *testing.T) {
	email := "email"
	password := "wrong-password"
	userID := "user-id"
	token := "generated-token"
	user := &users.User{
		ID:       userID,
		Password: &password,
	}

	generationMock := &domain.GenerationMock{}
	generationMock.On("ComparePassword", password, password).Return(true)
	generationMock.On("GenerateToken").Return(token, nil)

	userStorageMock := &users.StorageMock{}
	userStorageMock.On("GetByEmail", email).Return(user, nil)

	authorizationStorageMock := &StorageMock{}
	authorizationStorageMock.On("CreateAuthorization", userID, token).Return(nil)

	authorizationService := &AuthorizationService{
		AuthorizationStorage: authorizationStorageMock,
		UserStorage:          userStorageMock,
		Generation:           generationMock,
	}

	result, err := authorizationService.CreateAuthorization(email, password)
	expectedAuthorization := &Authorization{
		AccessToken: token,
		Owner: &users.User{
			ID: userID,
		},
	}

	generationMock.AssertCalled(t, "ComparePassword", password, password)
	generationMock.AssertCalled(t, "GenerateToken")
	userStorageMock.AssertCalled(t, "GetByEmail", email)
	authorizationStorageMock.AssertCalled(t, "CreateAuthorization", userID, token)
	assert.Equal(t, expectedAuthorization, result)
	assert.Empty(t, err)
}
