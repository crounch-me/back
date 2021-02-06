package utils

import (
	"github.com/crounch-me/back/internal/common/errors"
	"golang.org/x/crypto/bcrypt"
)

type Hash struct{}

func NewHash() HashLibrary {
	return &Hash{}
}

func (h Hash) Hash(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}

	return string(hash), nil
}

func (h Hash) Compare(s, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	return err == nil
}
