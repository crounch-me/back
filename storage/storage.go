package storage

import (
	"github.com/Sehsyha/crounch-back/model"
)

// Storage defines every data functions that we need
type Storage interface {
	CreateUser(user *model.User) error
	GetUserIDByToken(token string) (*string, error)
	GetUserByEmail(email string) (*model.User, error)

	CreateAuthorization(user *model.User) (*model.Authorization, error)

	CreateList(list *model.List) error
	GetOwnerLists(ownerID string) ([]*model.List, error)
}
