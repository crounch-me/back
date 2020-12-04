package listing

import (
	"time"

	"github.com/crounch-me/back/internal/common/errors"
)

// Storage defines every data functions that we need in lists /internal
type Storage interface {
	// Lists
	CreateList(id, name string, creationDate time.Time) *errors.Error
	GetUsersLists(userID string) ([]*List, *errors.Error)
	GetList(id string) (*List, *errors.Error)
	DeleteList(listID string) *errors.Error
	ArchiveList(listID string, archivationDate time.Time) *errors.Error

	// Products in list
	GetProductInList(productID string, listID string) (*ProductInListLink, *errors.Error)
	AddProductToList(productID string, listID string) *errors.Error
	DeleteProductsFromList(listID string) *errors.Error
	DeleteProductFromList(productID, listID string) *errors.Error
	GetProductsOfList(listID string) ([]*ProductInList, *errors.Error)
	UpdateProductInList(updateProductInList *UpdateProductInList, productID, listID string) (*ProductInListLink, *errors.Error)
}
