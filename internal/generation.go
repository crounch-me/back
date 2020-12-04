package internal

import "github.com/crounch-me/back/internal/common/errors"

type Generation interface {
	GenerateToken() (string, *errors.Error)
	GenerateID() (string, *errors.Error)
	HashPassword(string) (string, *errors.Error)
	ComparePassword(string, string) bool
}
