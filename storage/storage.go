package storage

import (
	"github.com/crounch-me/back/internal/account"
	"github.com/crounch-me/back/internal/contributors"
	"github.com/crounch-me/back/internal/list"
	"github.com/crounch-me/back/internal/products"
)

// Storage defines every data functions that we need
type Storage interface {
	account.Storage
	list.Storage
	products.Storage
	contributors.Storage
}
