package lists

import (
	"time"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/contributors"
	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
)

type ListService struct {
	ListStorage    Storage
  ProductStorage products.Storage
  ContributorStorage contributors.Storage
	Generation     domain.Generation
}

func (ls *ListService) CreateList(name, userID string) (*List, *domain.Error) {
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
		Contributors: []*users.User{
      {
        ID: userID,
      },
		},
		CreationDate: creationDate,
	}

	return list, nil
}

func (ls *ListService) GetUsersLists(userID string) ([]*List, *domain.Error) {
	lists, err := ls.ListStorage.GetUsersLists(userID)

	if err != nil {
		return nil, err
	}

	return lists, nil
}

func (ls *ListService) UpdateProductInList(updateProductInList *UpdateProductInList, productID, listID, userID string) (*ProductInListLink, *domain.Error) {
	_, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return nil, err
  }

  isUserAuthorized, err := ls.isUserAuthorized(listID, userID)
  if err != nil {
    return nil, err
  }

	if !isUserAuthorized {
		return nil, domain.NewError(domain.ForbiddenErrorCode)
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

func (ls *ListService) GetList(listID, userID string) (*List, *domain.Error) {
	list, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return nil, err
  }

  isUserAuthorized, err := ls.isUserAuthorized(listID, userID)
  if err != nil {
    return nil, err
  }

	if !isUserAuthorized {
		return nil, domain.NewError(domain.ForbiddenErrorCode)
	}

	list.Products, err = ls.ListStorage.GetProductsOfList(listID)
	if err != nil {
		return nil, err
	}

	return list, err
}

func (ls *ListService) DeleteProductFromList(productID, listID, userID string) *domain.Error {
	_, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return err
	}

  isUserAuthorized, err := ls.isUserAuthorized(listID, userID)
  if err != nil {
    return err
  }

	if !isUserAuthorized {
		return domain.NewError(domain.ForbiddenErrorCode)
	}

	_, err = ls.ProductStorage.GetProduct(productID)
	if err != nil {
		return err
	}

	return ls.ListStorage.DeleteProductFromList(productID, listID)
}

func (ls *ListService) AddProductToList(productID, listID, userID string) (*ProductInListLink, *domain.Error) {
	_, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return nil, err
	}

  isUserAuthorized, err := ls.isUserAuthorized(listID, userID)
  if err != nil {
    return nil, err
  }

	if !isUserAuthorized {
		return nil, domain.NewError(domain.ForbiddenErrorCode)
	}

	product, err := ls.ProductStorage.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	if !products.IsUserAuthorized(product, userID) {
		return nil, domain.NewError(domain.ForbiddenErrorCode)
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

	productInList = &ProductInListLink{
		ProductID: productID,
		ListID:    listID,
	}

	return productInList, nil
}

func (ls *ListService) DeleteList(listID, userID string) *domain.Error {
	_, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return err
  }

  isUserAuthorized, err := ls.isUserAuthorized(listID, userID)
  if err != nil {
    return  err
  }

	if !isUserAuthorized {
		return domain.NewError(domain.ForbiddenErrorCode)
	}

	err = ls.ListStorage.DeleteProductsFromList(listID)
	if err != nil {
		return err
	}

	return ls.ListStorage.DeleteList(listID)
}
