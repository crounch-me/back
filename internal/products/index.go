package products

import (
	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/user"
)

type ProductService struct {
	ProductStorage Storage
	Generation     internal.Generation
}

func (ps *ProductService) GetProduct(productID, userID string) (*Product, *internal.Error) {
	product, err := ps.ProductStorage.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	if !IsUserAuthorized(product, userID) {
		return nil, internal.NewError(internal.ForbiddenErrorCode)
	}

	return product, nil
}

func (ps *ProductService) CreateProduct(name, userID string) (*Product, *internal.Error) {
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
		Owner: &user.User{
			ID: userID,
		},
	}

	return product, nil
}

func (ps *ProductService) SearchDefaults(name string) ([]*Product, *internal.Error) {
	return ps.ProductStorage.SearchDefaults(name, user.AdminID)
}
