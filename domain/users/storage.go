package users

import "github.com/crounch-me/back/domain"

type Storage interface {
	CreateUser(id, email, password string) *domain.Error
	GetUserIDByToken(token string) (*string, *domain.Error)
	GetByEmail(email string) (*User, *domain.Error)
}
