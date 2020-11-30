package app

import (
	"errors"

	"github.com/crounch-me/back/internal/common/utils"
	"github.com/crounch-me/back/internal/user/domain/authorizations"
	"github.com/crounch-me/back/internal/user/domain/users"
)

type UserService struct {
	authorizationsRepository authorizations.Repository
	generationLibrary        utils.GenerationLibrary
	hashLibrary              utils.HashLibrary
	usersRepository          users.Repository
}

var (
	ErrUserNotFound = errors.New("user not found")

	NotFoundIndex = -1
)

func NewUserService(authorizationsRepository authorizations.Repository, usersRepository users.Repository) (*UserService, error) {
	if authorizationsRepository == nil {
		return nil, errors.New("authorizationsRepository is nil")
	}

	if usersRepository == nil {
		return nil, errors.New("usersRepository is nil")
	}

	return &UserService{
		authorizationsRepository: authorizationsRepository,
		usersRepository:          usersRepository,
	}, nil
}

func (s *UserService) Signup(email, password string) error {
	_, err := s.usersRepository.FindByEmail(email)
	if err != nil {
		if err != ErrUserNotFound {
			return err
		}
	} else {
		return errors.New("user already registered")
	}

	uuid, err := s.generationLibrary.UUID()
	if err != nil {
		return err
	}

	hashedPassword, err := s.hashLibrary.Hash(password)
	if err != nil {
		return err
	}

	u, err := users.NewUser(uuid, email, hashedPassword)
	if err != nil {
		return err
	}

	return s.usersRepository.AddUser(u)
}

func (s *UserService) Login(email, password string) (string, error) {
	u, err := s.usersRepository.FindByEmail(email)
	if err != nil {
		return "", err
	}

	isPasswordValid := s.hashLibrary.Compare(password, u.Password())
	if !isPasswordValid {
		return "", errors.New("invalid password")
	}

	token, err := s.generationLibrary.Token()
	if err != nil {
		return "", err
	}

	err = s.authorizationsRepository.AddAuthorization(u.UUID(), token)
	if err != nil {
		return "", err
	}

	return token, nil
}
