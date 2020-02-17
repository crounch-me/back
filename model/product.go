package model

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name" validate:"required,lt=61"`
	Owner *User  `json:"owner"`
}
