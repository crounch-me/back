package authorization

import (
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
)

type AuthorizationService struct {
	UserStorage          users.Storage
	AuthorizationStorage Storage
	Generation           domain.Generation
}

func (as *AuthorizationService) CreateAuthorization(email, password string) (*Authorization, *domain.Error) {
	user, err := as.UserStorage.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	isPasswordEqual := as.Generation.ComparePassword(*user.Password, password)
	if !isPasswordEqual {
		return nil, domain.NewError(WrongPasswordErrorCode)
	}

	token, err := as.Generation.GenerateToken()

	if err != nil {
		return nil, err
	}

	err = as.AuthorizationStorage.CreateAuthorization(user.ID, token)
	if err != nil {
		return nil, err
	}

	authorization := &Authorization{
		AccessToken: token,
		Owner: &users.User{
			ID: user.ID,
		},
	}

	return authorization, nil
}

func (as *AuthorizationService) Logout(userID, token string) *domain.Error {
  return as.AuthorizationStorage.DeleteAuthorization(userID, token)
}
