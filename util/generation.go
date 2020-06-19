package util

import (
	"math/rand"

	"github.com/crounch-me/back/domain"
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

func (g *GenerationImpl) GenerateToken() string {
	return RandString(tokenLength)
}

func (g *GenerationImpl) GenerateID() (string, *domain.Error) {
	id, err := uuid.NewV4()

	if err != nil {
		return "", domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}

	return id.String(), nil
}

func (g *GenerationImpl) HashPassword(password string) (string, *domain.Error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}
	return string(hashedPassword), nil
}

func (g *GenerationImpl) ComparePassword(hashedPassword, givenPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(givenPassword))
	return err == nil
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
