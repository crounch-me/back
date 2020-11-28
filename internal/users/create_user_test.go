package users

import (
	"testing"

	"github.com/crounch-me/back/internal"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserHashPasswordError(t *testing.T) {
	email := "email"
	password := "password"

	generationMock := &internal.GenerationMock{}
	generationMock.On("HashPassword", password).Return("", internal.NewError(internal.UnknownErrorCode))

	userService := &UserService{
		Generation: generationMock,
	}
	result, err := userService.CreateUser(email, password)

	assert.Empty(t, result)
	assert.Equal(t, internal.UnknownErrorCode, err.Code)
}

func TestCreateUserGenerateIDError(t *testing.T) {
	email := "email"
	password := "password"
	hashedPassword := "hashed-password"

	generationMock := &internal.GenerationMock{}
	generationMock.On("HashPassword", password).Return(hashedPassword, nil)
	generationMock.On("GenerateID").Return("", internal.NewError(internal.UnknownErrorCode))

	userService := &UserService{
		Generation: generationMock,
	}
	result, err := userService.CreateUser(email, password)

	assert.Empty(t, result)
	assert.Equal(t, internal.UnknownErrorCode, err.Code)
}

func TestCreateUserDuplicatedUserError(t *testing.T) {
	email := "email"
	password := "password"
	hashedPassword := "hashed-password"
	generatedID := "generated-id"

	foundUser := &User{}

	generationMock := &internal.GenerationMock{}
	generationMock.On("HashPassword", password).Return(hashedPassword, nil)
	generationMock.On("GenerateID").Return(generatedID, nil)

	userStorageMock := &StorageMock{}
	userStorageMock.On("GetByEmail", email).Return(foundUser, nil)

	userService := &UserService{
		Generation:  generationMock,
		UserStorage: userStorageMock,
	}
	result, err := userService.CreateUser(email, password)

	assert.Empty(t, result)
	assert.Equal(t, DuplicateUserErrorCode, err.Code)
}

func TestCreateUserGetByEmailError(t *testing.T) {
	email := "email"
	password := "password"
	hashedPassword := "hashed-password"
	generatedID := "generated-id"

	generationMock := &internal.GenerationMock{}
	generationMock.On("HashPassword", password).Return(hashedPassword, nil)
	generationMock.On("GenerateID").Return(generatedID, nil)

	userStorageMock := &StorageMock{}
	userStorageMock.On("GetByEmail", email).Return(nil, internal.NewError(internal.UnknownErrorCode))

	userService := &UserService{
		Generation:  generationMock,
		UserStorage: userStorageMock,
	}
	result, err := userService.CreateUser(email, password)

	assert.Empty(t, result)
	assert.Equal(t, internal.UnknownErrorCode, err.Code)
}

func TestCreateUserCreateUserError(t *testing.T) {
	email := "email"
	password := "password"
	hashedPassword := "hashed-password"
	generatedID := "generated-id"

	generationMock := &internal.GenerationMock{}
	generationMock.On("HashPassword", password).Return(hashedPassword, nil)
	generationMock.On("GenerateID").Return(generatedID, nil)

	userStorageMock := &StorageMock{}
	userStorageMock.On("GetByEmail", email).Return(nil, internal.NewError(UserNotFoundErrorCode))
	userStorageMock.On("CreateUser", generatedID, email, hashedPassword).Return(internal.NewError(internal.UnknownErrorCode))

	userService := &UserService{
		Generation:  generationMock,
		UserStorage: userStorageMock,
	}
	result, err := userService.CreateUser(email, password)

	assert.Empty(t, result)
	assert.Equal(t, internal.UnknownErrorCode, err.Code)
}

func TestCreateUserOK(t *testing.T) {
	email := "email"
	password := "password"
	hashedPassword := "hashed-password"
	generatedID := "generated-id"

	generationMock := &internal.GenerationMock{}
	generationMock.On("HashPassword", password).Return(hashedPassword, nil)
	generationMock.On("GenerateID").Return(generatedID, nil)

	userStorageMock := &StorageMock{}
	userStorageMock.On("GetByEmail", email).Return(nil, internal.NewError(UserNotFoundErrorCode))
	userStorageMock.On("CreateUser", generatedID, email, hashedPassword).Return(nil)

	userService := &UserService{
		Generation:  generationMock,
		UserStorage: userStorageMock,
	}
	result, err := userService.CreateUser(email, password)

	expectedUser := &User{
		ID:    generatedID,
		Email: email,
	}

	assert.Equal(t, expectedUser, result)
	assert.Empty(t, err)
}
