package products

import (
	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/account"
	"github.com/crounch-me/back/internal/common/errors"
)

type ProductService struct {
	ProductStorage Storage
	Generation     internal.Generation
}

func (ps *ProductService) GetProduct(productID, userID string) (*Product, *errors.Error) {
	product, err := ps.ProductStorage.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	if !IsUserAuthorized(product, userID) {
		return nil, errors.NewError(errors.ForbiddenErrorCode)
	}

	return product, nil
}

func (ps *ProductService) CreateProduct(name, userID string) (*Product, *errors.Error) {
	id, err := ps.Generation.GenerateID()
	if err != nil {
		return nil, err
	}

	err = ps.ProductStorage.CreateProduct(id, name, userID)
	if err != nil {
		return nil, err
	}

	product := &Product{
		ID:   id,
		Name: name,
		Owner: &account.User{
			ID: userID,
		},
	}

	return product, nil
}

func (ps *ProductService) SearchDefaults(name string) ([]*Product, *errors.Error) {
	return ps.ProductStorage.SearchDefaults(name, account.AdminID)
}
