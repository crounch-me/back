package authorization

import "github.com/crounch-me/back/domain/users"

type Authorization struct {
	AccessToken  string      `json:"accessToken"`
	ExpireDate   string      `json:"-"`
	RefreshToken string      `json:"-"`
	Owner        *users.User `json:"owner"`
}
