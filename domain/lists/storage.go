package lists

import "github.com/crounch-me/back/domain"

// ListStorage defines every data functions that we need
type Storage interface {
	CreateList(id, name, ownerID string) *domain.Error
	GetOwnersLists(ownerID string) ([]*List, *domain.Error)
	GetList(id string) (*List, *domain.Error)
	GetProductInList(productID string, listID string) (*ProductInList, *domain.Error)
	AddProductToList(productID string, listID string) *domain.Error
}
