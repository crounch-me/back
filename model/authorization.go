package model

type Authorization struct {
	AccessToken  string `json:"accessToken"`
	ExpireDate   string `json:"expireDate"`
	RefreshToken string `json:"refreshToken"`
}
