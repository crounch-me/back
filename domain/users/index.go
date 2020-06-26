package users

import (
	"github.com/crounch-me/back/domain"
)

const (
	AdminID = "00000000-0000-0000-0000-000000000000"
)

type UserService struct {
	UserStorage Storage
	Generation  domain.Generation
}

func (us *UserService) CreateUser(email string, password string) (*User, *domain.Error) {
	hashedPassword, err := us.Generation.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUserID, err := us.Generation.GenerateID()
	if err != nil {
		return nil, err
	}

	_, err = us.UserStorage.GetByEmail(email)
	if err == nil {
		return nil, domain.NewError(DuplicateUserErrorCode)
	} else if err.Code != UserNotFoundErrorCode {
		return nil, err
	}

	err = us.UserStorage.CreateUser(newUserID, email, hashedPassword)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:    newUserID,
		Email: email,
	}

	return user, nil
}

func (us *UserService) GetByEmail(email string) (*User, *domain.Error) {
	user, err := us.UserStorage.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
