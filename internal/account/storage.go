package account

import "github.com/crounch-me/back/internal/common/errors"

type Storage interface {
	CreateUser(id, email, password string) *errors.Error
	GetUserIDByToken(token string) (*string, *errors.Error)
	GetByEmail(email string) (*User, *errors.Error)
	GetByToken(token string) (*User, *errors.Error)
}
