package products

import "github.com/crounch-me/back/internal"

type Storage interface {
	CreateProduct(id, name, ownerID string) *internal.Error
	GetProduct(id string) (*Product, *internal.Error)
	SearchDefaults(name string, id string) ([]*Product, *internal.Error)
}
