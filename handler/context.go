package handler

import (
	"github.com/Sehsyha/crounch-back/configuration"
	"github.com/Sehsyha/crounch-back/storage"
	"github.com/Sehsyha/crounch-back/storage/mock"
	"github.com/Sehsyha/crounch-back/storage/neo"
)

type Context struct {
	Storage storage.Storage
	Config  *configuration.Config
}

func NewContext(config *configuration.Config) *Context {
	var storage storage.Storage

	if config.Mock {
		storage = mock.NewStorageMock()
	} else {
		storage = neo.NewNeoStorage()
	}

	return &Context{
		Storage: storage,
		Config:  config,
	}
}
