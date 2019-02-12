package handler

import (
	"github.com/Sehsyha/crounch-back/config"
	"github.com/Sehsyha/crounch-back/storage"
	"github.com/Sehsyha/crounch-back/storage/mock"
	"github.com/Sehsyha/crounch-back/storage/neo"
)

type Context struct {
	Storage storage.Storage
	Config  *config.Config
}

func NewContext(config *config.Config) *Context {
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
