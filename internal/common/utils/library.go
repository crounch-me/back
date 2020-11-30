package utils

type GenerationLibrary interface {
	UUID() (string, error)
	Token() (string, error)
}

type HashLibrary interface {
	Hash(s string) (string, error)
	Compare(s, hash string) bool
}
