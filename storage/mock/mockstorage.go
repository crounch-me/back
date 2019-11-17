package mock

import (
	"github.com/Sehsyha/crounch-back/storage"

	testmock "github.com/stretchr/testify/mock"
)

type StorageMock struct {
	testmock.Mock
}

func NewStorageMock() storage.Storage {
	return &StorageMock{}
}
