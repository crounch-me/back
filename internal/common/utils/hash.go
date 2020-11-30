package utils

import (
	"github.com/crounch-me/back/internal"
	"golang.org/x/crypto/bcrypt"
)

func Hash(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", internal.NewError(internal.UnknownErrorCode).WithCause(err)
	}
	return string(hash), nil
}

func CompareWithHash(s, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	return err == nil
}
