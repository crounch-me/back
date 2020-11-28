package contributors

import "github.com/crounch-me/back/internal"

// Storage defines every data functions that we need in contributors /internal
type Storage interface {
	CreateContributor(listID, userID string) *internal.Error
	GetContributorsIDs(listID string) ([]string, *internal.Error)
}
