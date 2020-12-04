package account

import (
	"github.com/crounch-me/back/internal"
)

const (
	AdminID = "00000000-0000-0000-0000-000000000000"
)

type UserService struct {
	UserStorage Storage
	Generation  internal.Generation
}
