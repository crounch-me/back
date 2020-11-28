package internal

import (
	"github.com/stretchr/testify/mock"
)

type GenerationMock struct {
	mock.Mock
}

func (g *GenerationMock) GenerateID() (string, *Error) {
	args := g.Called()
	err := args.Error(1)
	if err == nil {
		return args.String(0), nil
	}

	return args.String(0), err.(*Error)
}

func (g *GenerationMock) GenerateToken() (string, *Error) {
	args := g.Called()
	err := args.Error(1)
	if err == nil {
		return args.String(0), nil
	}
	return args.String(0), err.(*Error)
}

func (g *GenerationMock) HashPassword(password string) (string, *Error) {
	args := g.Called(password)
	err := args.Error(1)
	if err == nil {
		return args.String(0), nil
	}
	return "", err.(*Error)
}

func (g *GenerationMock) ComparePassword(hashedPassword, givenPassword string) bool {
	args := g.Called(hashedPassword, givenPassword)
	return args.Bool(0)
}
