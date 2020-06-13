package lists

import (
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
	"github.com/crounch-me/back/util"
)

type ListService struct {
	ListStorage    ListStorage
	ProductStorage products.ProductStorage
}

func (ls *ListService) CreateList(name, userID string) (*List, *domain.Error) {
	id, err := util.GenerateID()
	if err != nil {
		return nil, err
	}

	list := &List{
		ID:   id,
		Name: name,
		Owner: &users.User{
			ID: userID,
		},
	}

	err = ls.ListStorage.CreateList(list)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (ls *ListService) GetOwnerLists(ownerID string) ([]*List, *domain.Error) {
	lists, err := ls.ListStorage.GetOwnerLists(ownerID)

	if err != nil {
		return nil, err
	}

	return lists, nil
}

func (ls *ListService) GetList(listID, userID string) (*List, *domain.Error) {
	list, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return nil, err
	}

	if !IsUserAuthorized(list, userID) {
		return nil, domain.NewError(domain.UnauthorizedErrorCode)
	}

	return list, err
}

func (ls *ListService) AddProductToList(productID, listID, userID string) (*ProductInList, *domain.Error) {
	list, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return nil, err
	}

	if !IsUserAuthorized(list, userID) {
		return nil, domain.NewError(domain.UnauthorizedErrorCode)
	}

	product, err := ls.ProductStorage.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	if !products.IsUserAuthorized(product, userID) {
		return nil, domain.NewError(domain.UnauthorizedErrorCode)
	}

	productInList, err := ls.ListStorage.GetProductInList(productID, listID)
	if err == nil {
		return nil, domain.NewError(DuplicateProductInListErrorCode)
	} else if err.Code != ProductInListNotFoundErrorCode {
		return nil, err
	}

	err = ls.ListStorage.AddProductToList(productID, listID)
	if err != nil {
		return nil, err
	}

	productInList = &ProductInList{
		ProductID: productID,
		ListID:    listID,
	}

	return productInList, nil
}
