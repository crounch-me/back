package users

import (
	"testing"

	"github.com/crounch-me/back/internal"
	"github.com/stretchr/testify/assert"
)

func TestGetByTokenGetByTokenError(t *testing.T) {
	token := "token"

	userStorageMock := &StorageMock{}
	userStorageMock.On("GetByToken", token).Return(nil, internal.NewError(UserNotFoundErrorCode))

	userService := &UserService{
		UserStorage: userStorageMock,
	}
	result, err := userService.GetByToken(token)

	userStorageMock.AssertCalled(t, "GetByToken", token)

	assert.Empty(t, result)
	assert.Equal(t, UserNotFoundErrorCode, err.Code)
}

func TestGetByTokenOK(t *testing.T) {
	token := "token"

	user := &User{
		ID:    "user-id",
		Email: "user-email",
	}

	userStorageMock := &StorageMock{}
	userStorageMock.On("GetByToken", token).Return(user, nil)

	userService := &UserService{
		UserStorage: userStorageMock,
	}
	result, err := userService.GetByToken(token)

	userStorageMock.AssertCalled(t, "GetByToken", token)

	assert.Equal(t, user, result)
	assert.Empty(t, err)
}
