package contributors

import (
	"github.com/crounch-me/back/domain"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (sm *StorageMock) CreateContributor(listID, userID string) *domain.Error {
	args := sm.Called(listID, userID)
	err := args.Error(0)

  if err == nil {
		return nil
  }

	return err.(*domain.Error)
}

func (sm *StorageMock) GetContributorsIDs(listID string) ([]string, *domain.Error)  {
  args := sm.Called(listID)
  err := args.Error(1)

  if err != nil {
    return []string{}, err.(*domain.Error)
  }

  return args.Get(0).([]string), nil
}
