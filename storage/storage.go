package storage

import "github.com/Sehsyha/crounch-back/model"

type Storage interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)

	CreateAuthorization(user *model.User) (*model.Authorization, error)
}
