package app

import (
	"errors"

	"github.com/crounch-me/back/internal/account/domain/authorizations"
	"github.com/crounch-me/back/internal/account/domain/users"
	"github.com/crounch-me/back/internal/common/utils"
)

type AccountService struct {
	authorizationsRepository authorizations.Repository
	generationLibrary        utils.GenerationLibrary
	hashLibrary              utils.HashLibrary
	usersRepository          users.Repository
}

const (
	NotFoundIndex = -1
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func NewAccountService(
	authorizationsRepository authorizations.Repository,
	generationLibrary utils.GenerationLibrary,
	hashLibrary utils.HashLibrary,
	usersRepository users.Repository,
) (*AccountService, error) {
	if authorizationsRepository == nil {
		return nil, errors.New("authorizationsRepository is nil")
	}

	if generationLibrary == nil {
		return nil, errors.New("generationLibrary is nil")
	}

	if hashLibrary == nil {
		return nil, errors.New("hashLibrary is nil")
	}

	if usersRepository == nil {
		return nil, errors.New("usersRepository is nil")
	}

	return &AccountService{
		authorizationsRepository: authorizationsRepository,
		generationLibrary:        generationLibrary,
		hashLibrary:              hashLibrary,
		usersRepository:          usersRepository,
	}, nil
}

func (s *AccountService) Signup(email, password string) error {
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

func (s *AccountService) Login(email, password string) (string, error) {
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

func (s *AccountService) Logout(userUUID, token string) error {
	return s.authorizationsRepository.RemoveAuthorization(userUUID, token)
}

func (s *AccountService) GetUserUUIDByToken(token string) (string, error) {
	return s.authorizationsRepository.GetUserUUIDByToken(token)
}
