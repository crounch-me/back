package users

import (
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/util"
)

type UserService struct {
	UserStorage UserStorage
}

func (us *UserService) CreateUser(email string, password string) (*User, *domain.Error) {
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUserID, err := util.GenerateID()
	if err != nil {
		return nil, err
	}

	_, err = us.UserStorage.GetUserByEmail(email)
	if err == nil {
		return nil, domain.NewError(DuplicateUserErrorCode)
	} else if err.Code != UserNotFoundErrorCode {
		return nil, err
	}

	user := &User{
		ID:       newUserID,
		Email:    email,
		Password: &hashedPassword,
	}

	err = us.UserStorage.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) GetUserByEmail(email string) (*User, *domain.Error) {
	user, err := us.UserStorage.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
