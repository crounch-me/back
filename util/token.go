package util

import (
	"math/rand"
)

const (
	tokenLength = 42
	charset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func GenerateToken() string {
	return RandStringRunes(tokenLength)
}

func RandStringRunes(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
