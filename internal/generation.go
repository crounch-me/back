package internal

type Generation interface {
	GenerateToken() (string, *Error)
	GenerateID() (string, *Error)
	HashPassword(string) (string, *Error)
	ComparePassword(string, string) bool
}
