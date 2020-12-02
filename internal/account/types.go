package account

// User represents an user of the application
type User struct {
	ID       string  `json:"id"`
	Email    string  `json:"email,omitempty" validate:"required,email"`
	Password *string `json:"password,omitempty" validate:"required,gt=3"`
}
