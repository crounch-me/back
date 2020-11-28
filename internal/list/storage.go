package list

import (
	"time"

	"github.com/crounch-me/back/internal"
)

// Storage defines every data functions that we need in lists /internal
type Storage interface {
	// Lists
	CreateList(id, name string, creationDate time.Time) *internal.Error
	GetUsersLists(userID string) ([]*List, *internal.Error)
	GetList(id string) (*List, *internal.Error)
	DeleteList(listID string) *internal.Error
	ArchiveList(listID string, archivationDate time.Time) *internal.Error

	// Products in list
	GetProductInList(productID string, listID string) (*ProductInListLink, *internal.Error)
	AddProductToList(productID string, listID string) *internal.Error
	DeleteProductsFromList(listID string) *internal.Error
	DeleteProductFromList(productID, listID string) *internal.Error
	GetProductsOfList(listID string) ([]*ProductInList, *internal.Error)
	UpdateProductInList(updateProductInList *UpdateProductInList, productID, listID string) (*ProductInListLink, *internal.Error)
}
