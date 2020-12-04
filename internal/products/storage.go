package products

import "github.com/crounch-me/back/internal/common/errors"

type Storage interface {
	CreateProduct(id, name, ownerID string) *errors.Error
	GetProduct(id string) (*Product, *errors.Error)
	SearchDefaults(name string, id string) ([]*Product, *errors.Error)
}
