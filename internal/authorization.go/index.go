package authorization

import (
	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/users"
)

type AuthorizationService struct {
	UserStorage          users.Storage
	AuthorizationStorage Storage
	Generation           internal.Generation
}

func (as *AuthorizationService) CreateAuthorization(email, password string) (*Authorization, *internal.Error) {
	user, err := as.UserStorage.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	isPasswordEqual := as.Generation.ComparePassword(*user.Password, password)
	if !isPasswordEqual {
		return nil, internal.NewError(WrongPasswordErrorCode)
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

func (as *AuthorizationService) Logout(token string) *internal.Error {
	user, err := as.UserStorage.GetByToken(token)

	if err != nil {
		return err
	}

	return as.AuthorizationStorage.DeleteAuthorization(user.ID, token)
}
