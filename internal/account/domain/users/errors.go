package users

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")

	ErrDuplicateUser = errors.New("user already exists")
)
