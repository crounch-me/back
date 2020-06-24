package products

import "github.com/crounch-me/back/domain"

type Storage interface {
	CreateProduct(product *Product) *domain.Error
	GetProduct(id string) (*Product, *domain.Error)
}
