package authorization

import "github.com/crounch-me/back/internal/user"

type Authorization struct {
	AccessToken  string     `json:"accessToken"`
	ExpireDate   string     `json:"-"`
	RefreshToken string     `json:"-"`
	Owner        *user.User `json:"owner"`
}
