package users

type Repository interface {
	AddUser(user *User) error
	FindByEmail(email string) (*User, error)
	FindByToken(token string) (*User, error)
}
