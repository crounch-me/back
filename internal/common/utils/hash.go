package utils

import (
	"github.com/crounch-me/back/internal/common/errors"
	"golang.org/x/crypto/bcrypt"
)

func Hash(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}
	return string(hash), nil
}

func CompareWithHash(s, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	return err == nil
}
