package util

import (
	"math/rand"

	"github.com/crounch-me/back/internal/common/errors"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenLength   = 42
	charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type GenerationImpl struct{}

func (g *GenerationImpl) GenerateToken() (string, *errors.Error) {
	return g.GenerateID()
}

func (g *GenerationImpl) GenerateID() (string, *errors.Error) {
	id, err := GenerateID()

	if err != nil {
		return "", errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}

	return id, nil
}

func (g *GenerationImpl) HashPassword(password string) (string, *errors.Error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}
	return string(hashedPassword), nil
}

func (g *GenerationImpl) ComparePassword(hashedPassword, givenPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(givenPassword))
	return err == nil
}

func GenerateID() (string, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func RandString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(charset) {
			b[i] = charset[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
