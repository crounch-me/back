package model

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner *User  `json:"owner"`
}
