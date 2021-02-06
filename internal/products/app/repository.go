package app

import "github.com/crounch-me/back/internal/products/domain/products"

type Repository interface {
	SaveProduct(p *products.Product) error
	ReadByID(uuid string) (*products.Product, error)
	SearchDefaults(name string) ([]*products.Product, error)
}
