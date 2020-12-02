package account

import "github.com/crounch-me/back/internal"

type Storage interface {
	CreateUser(id, email, password string) *internal.Error
	GetUserIDByToken(token string) (*string, *internal.Error)
	GetByEmail(email string) (*User, *internal.Error)
	GetByToken(token string) (*User, *internal.Error)
}
