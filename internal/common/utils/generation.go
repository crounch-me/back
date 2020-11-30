package utils

import "github.com/gofrs/uuid"

type Generation struct{}

func (g Generation) UUID() (string, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (g Generation) Token() (string, error) {
	return g.Token()
}
