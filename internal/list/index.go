package list

import (
	"time"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/contributors"
	"github.com/crounch-me/back/internal/products"
	"github.com/crounch-me/back/internal/user"
)

type ListService struct {
	ListStorage        Storage
	ProductStorage     products.Storage
	ContributorStorage contributors.Storage
	Generation         internal.Generation
}

func (ls *ListService) CreateList(name, userID string) (*List, *internal.Error) {
	id, err := ls.Generation.GenerateID()
	if err != nil {
		return nil, err
	}

	creationDate := time.Now()
	err = ls.ListStorage.CreateList(id, name, creationDate)

	if err != nil {
		return nil, err
	}

	err = ls.ContributorStorage.CreateContributor(id, userID)
	if err != nil {
		return nil, err
	}

	list := &List{
		ID:   id,
		Name: name,
		Contributors: []*user.User{
			{
				ID: userID,
			},
		},
		CreationDate: creationDate,
	}

	return list, nil
}

func (ls *ListService) GetUsersLists(userID string) ([]*List, *internal.Error) {
	lists, err := ls.ListStorage.GetUsersLists(userID)

	if err != nil {
		return nil, err
	}

	return lists, nil
}

func (ls *ListService) ArchiveList(listID, userID string) (*List, *internal.Error) {
	isUserAuthorized, err := ls.isUserAuthorized(listID, userID)
	if err != nil {
		return nil, err
	}

	if !isUserAuthorized {
		return nil, internal.NewError(internal.ForbiddenErrorCode)
	}

	list, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return nil, err
	}

	archivationDate := time.Now()

	err = ls.ListStorage.ArchiveList(listID, archivationDate)
	if err != nil {
		return nil, err
	}

	list.ArchivationDate = &archivationDate

	return list, nil
}

func (ls *ListService) UpdateProductInList(updateProductInList *UpdateProductInList, productID, listID, userID string) (*ProductInListLink, *internal.Error) {
	_, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return nil, err
	}

	isUserAuthorized, err := ls.isUserAuthorized(listID, userID)
	if err != nil {
		return nil, err
	}

	if !isUserAuthorized {
		return nil, internal.NewError(internal.ForbiddenErrorCode)
	}

	_, err = ls.ProductStorage.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	productInListLink, err := ls.ListStorage.UpdateProductInList(updateProductInList, productID, listID)
	if err != nil {
		return nil, err
	}

	return productInListLink, nil
}

func (ls *ListService) GetList(listID, userID string) (*List, *internal.Error) {
	list, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return nil, err
	}

	isUserAuthorized, err := ls.isUserAuthorized(listID, userID)
	if err != nil {
		return nil, err
	}

	if !isUserAuthorized {
		return nil, internal.NewError(internal.ForbiddenErrorCode)
	}

	list.Products, err = ls.ListStorage.GetProductsOfList(listID)
	if err != nil {
		return nil, err
	}

	return list, err
}

func (ls *ListService) DeleteProductFromList(productID, listID, userID string) *internal.Error {
	_, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return err
	}

	isUserAuthorized, err := ls.isUserAuthorized(listID, userID)
	if err != nil {
		return err
	}

	if !isUserAuthorized {
		return internal.NewError(internal.ForbiddenErrorCode)
	}

	_, err = ls.ProductStorage.GetProduct(productID)
	if err != nil {
		return err
	}

	return ls.ListStorage.DeleteProductFromList(productID, listID)
}

func (ls *ListService) AddProductToList(productID, listID, userID string) (*ProductInListLink, *internal.Error) {
	_, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return nil, err
	}

	isUserAuthorized, err := ls.isUserAuthorized(listID, userID)
	if err != nil {
		return nil, err
	}

	if !isUserAuthorized {
		return nil, internal.NewError(internal.ForbiddenErrorCode)
	}

	product, err := ls.ProductStorage.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	if !products.IsUserAuthorized(product, userID) {
		return nil, internal.NewError(internal.ForbiddenErrorCode)
	}

	productInList, err := ls.ListStorage.GetProductInList(productID, listID)
	if err == nil {
		return nil, internal.NewError(DuplicateProductInListErrorCode)
	} else if err.Code != ProductInListNotFoundErrorCode {
		return nil, err
	}

	err = ls.ListStorage.AddProductToList(productID, listID)
	if err != nil {
		return nil, err
	}

	productInList = &ProductInListLink{
		ProductID: productID,
		ListID:    listID,
	}

	return productInList, nil
}

func (ls *ListService) DeleteList(listID, userID string) *internal.Error {
	_, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return err
	}

	isUserAuthorized, err := ls.isUserAuthorized(listID, userID)
	if err != nil {
		return err
	}

	if !isUserAuthorized {
		return internal.NewError(internal.ForbiddenErrorCode)
	}

	err = ls.ListStorage.DeleteProductsFromList(listID)
	if err != nil {
		return err
	}

	return ls.ListStorage.DeleteList(listID)
}
