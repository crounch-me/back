package domain

type Generation interface {
	GenerateToken() string
	GenerateID() (string, *Error)
	HashPassword(string) (string, *Error)
	ComparePassword(string, string) bool
}
