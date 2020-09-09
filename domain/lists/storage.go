package lists

import (
	"time"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/products"
)

// ListStorage defines every data functions that we need
type Storage interface {
	// Lists
	CreateList(id, name, ownerID string, creationDate time.Time) *domain.Error
	GetOwnersLists(ownerID string) ([]*List, *domain.Error)
	GetList(id string) (*List, *domain.Error)
	DeleteList(listID string) *domain.Error

	// Products in list
	GetProductInList(productID string, listID string) (*ProductInList, *domain.Error)
	AddProductToList(productID string, listID string) *domain.Error
	DeleteProductsFromList(listID string) *domain.Error
	DeleteProductFromList(productID, listID string) *domain.Error
	GetProductsOfList(listID string) ([]*products.Product, *domain.Error)
}
