package users

import "github.com/crounch-me/back/domain"

type UserStorage interface {
	CreateUser(user *User) *domain.Error
	GetUserIDByToken(token string) (*string, *domain.Error)
	GetUserByEmail(email string) (*User, *domain.Error)
}
