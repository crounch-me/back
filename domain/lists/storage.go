package lists

import (
	"time"

	"github.com/crounch-me/back/domain"
)

// Storage defines every data functions that we need in lists domain
type Storage interface {
	// Lists
	CreateList(id, name string, creationDate time.Time) *domain.Error
	GetUsersLists(userID string) ([]*List, *domain.Error)
	GetList(id string) (*List, *domain.Error)
	DeleteList(listID string) *domain.Error
	ArchiveList(listID string, archivationDate time.Time) *domain.Error

	// Products in list
	GetProductInList(productID string, listID string) (*ProductInListLink, *domain.Error)
	AddProductToList(productID string, listID string) *domain.Error
	DeleteProductsFromList(listID string) *domain.Error
	DeleteProductFromList(productID, listID string) *domain.Error
	GetProductsOfList(listID string) ([]*ProductInList, *domain.Error)
	UpdateProductInList(updateProductInList *UpdateProductInList, productID, listID string) (*ProductInListLink, *domain.Error)
}
