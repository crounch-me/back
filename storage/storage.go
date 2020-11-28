package storage

import (
	"github.com/crounch-me/back/internal/authorization.go"
	"github.com/crounch-me/back/internal/contributors"
	"github.com/crounch-me/back/internal/list"
	"github.com/crounch-me/back/internal/products"
	"github.com/crounch-me/back/internal/users"
)

// Storage defines every data functions that we need
type Storage interface {
	users.Storage
	authorization.Storage
	list.Storage
	products.Storage
	contributors.Storage
}
