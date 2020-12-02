package users

import "errors"

type User struct {
	uuid     string
	email    string
	password string
}

func NewUser(uuid, email, password string) (*User, error) {
	if uuid == "" {
		return nil, errors.New("empty product uuid")
	}

	if email == "" {
		return nil, errors.New("empty product email")
	}

	if password == "" {
		return nil, errors.New("empty product password")
	}

	return &User{
		uuid:     uuid,
		email:    email,
		password: password,
	}, nil
}

func (u User) UUID() string {
	return u.uuid
}

func (u User) Password() string {
	return u.password
}
