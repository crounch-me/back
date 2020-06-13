package authorization

import (
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
	"github.com/crounch-me/back/util"
)

type AuthorizationService struct {
	UserStorage          users.UserStorage
	AuthorizationStorage AuthorizationStorage
}

func (as *AuthorizationService) CreateAuthorization(email, password string) (*Authorization, *domain.Error) {
	user, err := as.UserStorage.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	isPasswordEqual := util.ComparePassword(*user.Password, password)
	if isPasswordEqual {
		return nil, domain.NewError(WrongPasswordErrorCode)
	}

	token := util.GenerateToken()

	authorization, err := as.AuthorizationStorage.CreateAuthorization(user.ID, token)
	if err != nil {
		return nil, err
	}

	return authorization, nil
}
