package authorization

import (
	"testing"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/users"
	"github.com/stretchr/testify/assert"
)

func TestLogoutDeleteGetByTokenError(t *testing.T) {
	userID := "user-id"
	token := "token"

	userStorageMock := &users.StorageMock{}
	userStorageMock.On("GetByToken", token).Return(nil, internal.NewError(internal.UnknownErrorCode))

	authorizationStorageMock := &StorageMock{}
	authorizationStorageMock.On("DeleteAuthorization", userID, token).Return(internal.NewError(internal.UnknownErrorCode))

	authorizationService := &AuthorizationService{
		AuthorizationStorage: authorizationStorageMock,
		UserStorage:          userStorageMock,
	}

	err := authorizationService.Logout(token)

	userStorageMock.AssertCalled(t, "GetByToken", token)

	authorizationStorageMock.AssertNotCalled(t, "DeleteAuthorization")

	assert.Equal(t, internal.UnknownErrorCode, err.Code)
}
func TestLogoutDeleteAuthorizationError(t *testing.T) {
	userID := "user-id"
	token := "token"

	user := &users.User{
		ID: userID,
	}

	userStorageMock := &users.StorageMock{}
	userStorageMock.On("GetByToken", token).Return(user, nil)

	authorizationStorageMock := &StorageMock{}
	authorizationStorageMock.On("DeleteAuthorization", userID, token).Return(internal.NewError(internal.UnknownErrorCode))

	authorizationService := &AuthorizationService{
		AuthorizationStorage: authorizationStorageMock,
		UserStorage:          userStorageMock,
	}

	err := authorizationService.Logout(token)

	userStorageMock.AssertCalled(t, "GetByToken", token)
	authorizationStorageMock.AssertCalled(t, "DeleteAuthorization", userID, token)

	assert.Equal(t, internal.UnknownErrorCode, err.Code)
}

func TestLogoutDeleteAuthorizationOK(t *testing.T) {
	userID := "user-id"
	token := "token"

	user := &users.User{
		ID: userID,
	}

	userStorageMock := &users.StorageMock{}
	userStorageMock.On("GetByToken", token).Return(user, nil)

	authorizationStorageMock := &StorageMock{}
	authorizationStorageMock.On("DeleteAuthorization", userID, token).Return(nil)

	authorizationService := &AuthorizationService{
		AuthorizationStorage: authorizationStorageMock,
		UserStorage:          userStorageMock,
	}

	err := authorizationService.Logout(token)

	userStorageMock.AssertCalled(t, "GetByToken", token)
	authorizationStorageMock.AssertCalled(t, "DeleteAuthorization", userID, token)

	assert.Empty(t, err)
}
