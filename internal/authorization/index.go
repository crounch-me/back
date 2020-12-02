package authorization

import (
	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/account"
)

type AuthorizationService struct {
	UserStorage          account.Storage
	AuthorizationStorage Storage
	Generation           internal.Generation
}

func (as *AuthorizationService) CreateAuthorization(email, password string) (*Authorization, *internal.Error) {
	u, err := as.UserStorage.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	isPasswordEqual := as.Generation.ComparePassword(*u.Password, password)
	if !isPasswordEqual {
		return nil, internal.NewError(WrongPasswordErrorCode)
	}

	token, err := as.Generation.GenerateToken()

	if err != nil {
		return nil, err
	}

	err = as.AuthorizationStorage.CreateAuthorization(u.ID, token)
	if err != nil {
		return nil, err
	}

	authorization := &Authorization{
		AccessToken: token,
		Owner: &account.User{
			ID: u.ID,
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
