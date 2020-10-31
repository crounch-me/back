package authorization

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/stretchr/testify/assert"
)

func TestLogoutDeleteAuthorizationError(t *testing.T) {
	userID := "user-id"
  token := "token"

  authorizationStorageMock := &StorageMock{}
  authorizationStorageMock.On("DeleteAuthorization", userID, token).Return(domain.NewError(domain.UnknownErrorCode))

	authorizationService := &AuthorizationService{
    AuthorizationStorage: authorizationStorageMock,
	}

	err := authorizationService.Logout(token)

	authorizationStorageMock.AssertCalled(t, "DeleteAuthorization", userID, token)

	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestLogoutDeleteAuthorizationOK(t *testing.T) {
	userID := "user-id"
  token := "token"

  authorizationStorageMock := &StorageMock{}
  authorizationStorageMock.On("DeleteAuthorization", userID, token).Return(nil)

	authorizationService := &AuthorizationService{
    AuthorizationStorage: authorizationStorageMock,
	}

	err := authorizationService.Logout(token)

	authorizationStorageMock.AssertCalled(t, "DeleteAuthorization", userID, token)

	assert.Empty(t, err)
}
