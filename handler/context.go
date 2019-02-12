package handler

import (
	"github.com/Sehsyha/crounch-back/storage"
	"github.com/Sehsyha/crounch-back/storage/neo"
)

type Context struct {
	Storage storage.Storage
}

func NewContext() *Context {
	storage := neo.NewNeoStorage()
	return &Context{
		Storage: storage,
	}
}
