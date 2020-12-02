package authorization

import "github.com/crounch-me/back/internal/account"

type Authorization struct {
	AccessToken  string        `json:"accessToken"`
	ExpireDate   string        `json:"-"`
	RefreshToken string        `json:"-"`
	Owner        *account.User `json:"owner"`
}
