package authorization

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
	"github.com/stretchr/testify/assert"
)

func TestLogoutDeleteGetByTokenError(t *testing.T) {
	userID := "user-id"
  token := "token"

  userStorageMock := &users.StorageMock{}
  userStorageMock.On("GetByToken", token).Return(nil, domain.NewError(domain.UnknownErrorCode))

  authorizationStorageMock := &StorageMock{}
  authorizationStorageMock.On("DeleteAuthorization", userID, token).Return(domain.NewError(domain.UnknownErrorCode))

	authorizationService := &AuthorizationService{
    AuthorizationStorage: authorizationStorageMock,
    UserStorage: userStorageMock,
	}

	err := authorizationService.Logout(token)

  userStorageMock.AssertCalled(t, "GetByToken", token)

	authorizationStorageMock.AssertNotCalled(t, "DeleteAuthorization")

	assert.Equal(t, domain.UnknownErrorCode, err.Code)
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
  authorizationStorageMock.On("DeleteAuthorization", userID, token).Return(domain.NewError(domain.UnknownErrorCode))

	authorizationService := &AuthorizationService{
    AuthorizationStorage: authorizationStorageMock,
    UserStorage: userStorageMock,
	}

	err := authorizationService.Logout(token)

  userStorageMock.AssertCalled(t, "GetByToken", token)
	authorizationStorageMock.AssertCalled(t, "DeleteAuthorization", userID, token)

	assert.Equal(t, domain.UnknownErrorCode, err.Code)
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
    UserStorage: userStorageMock,
	}

	err := authorizationService.Logout(token)

  userStorageMock.AssertCalled(t, "GetByToken", token)
	authorizationStorageMock.AssertCalled(t, "DeleteAuthorization", userID, token)

	assert.Empty(t, err)
}
