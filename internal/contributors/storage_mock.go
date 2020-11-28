package contributors

import (
	"github.com/crounch-me/back/internal"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (sm *StorageMock) CreateContributor(listID, userID string) *internal.Error {
	args := sm.Called(listID, userID)
	err := args.Error(0)

	if err == nil {
		return nil
	}

	return err.(*internal.Error)
}

func (sm *StorageMock) GetContributorsIDs(listID string) ([]string, *internal.Error) {
	args := sm.Called(listID)
	err := args.Error(1)

	if err != nil {
		return []string{}, err.(*internal.Error)
	}

	return args.Get(0).([]string), nil
}
