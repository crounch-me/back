package contributors

import "github.com/crounch-me/back/internal/common/errors"

// Storage defines every data functions that we need in contributors /internal
type Storage interface {
	CreateContributor(listID, userID string) *errors.Error
	GetContributorsIDs(listID string) ([]string, *errors.Error)
}
