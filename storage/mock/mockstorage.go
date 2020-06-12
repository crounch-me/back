package mock

import (
	"github.com/crounch-me/back/storage"

	testmock "github.com/stretchr/testify/mock"
)

type StorageMock struct {
	testmock.Mock
}

func NewStorageMock() storage.Storage {
	return &StorageMock{}
}
