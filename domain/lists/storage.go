package lists

import "github.com/crounch-me/back/domain"

// ListStorage defines every data functions that we need
type ListStorage interface {
	CreateList(list *List) *domain.Error
	GetOwnerLists(ownerID string) ([]*List, *domain.Error)
	GetList(id string) (*List, *domain.Error)
	GetProductInList(productID string, listID string) (*ProductInList, *domain.Error)
	AddProductToList(productID string, listID string) *domain.Error
}
