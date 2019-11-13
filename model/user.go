package model

type User struct {
	ID       string  `json:"id"`
	Email    string  `json:"email" validate:"required,email"`
	Password *string `json:"password,omitempty" validate:"required,gt=3"`
}
