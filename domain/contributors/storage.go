package contributors

import (
	"github.com/crounch-me/back/domain"
)

// Storage defines every data functions that we need in contributors domain
type Storage interface {
  CreateContributor(listID, userID string) *domain.Error
  GetContributorsIDs(listID string) ([]string, *domain.Error)
}
