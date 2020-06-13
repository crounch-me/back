package products

import "github.com/crounch-me/back/domain"

type ProductStorage interface {
	CreateProduct(product *Product) *domain.Error
	GetProduct(id string) (*Product, *domain.Error)
}
