package utils

import "github.com/gofrs/uuid"

type generation struct{}

func NewGeneration() GenerationLibrary {
	return &generation{}
}

func (g generation) UUID() (string, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (g generation) Token() (string, error) {
	return g.UUID()
}
