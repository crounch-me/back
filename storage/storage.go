package storage

import "github.com/Sehsyha/crounch-back/model"

type Storage interface {
	CreateUser(u *model.User) error
}
