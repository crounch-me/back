package authorization

import "github.com/crounch-me/back/domain/users"

type Authorization struct {
	AccessToken  string      `json:"accessToken"`
	ExpireDate   string      `json:"expireDate"`
	RefreshToken string      `json:"refreshToken"`
	Owner        *users.User `json:"owner"`
}
