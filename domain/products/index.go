package products

import (
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
	"github.com/crounch-me/back/util"
)

type ProductService struct {
	ProductStorage ProductStorage
}

func (ps *ProductService) GetProduct(productID, userID string) (*Product, *domain.Error) {
	product, err := ps.ProductStorage.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	if !IsUserAuthorized(product, userID) {
		return nil, domain.NewError(domain.UnauthorizedErrorCode)
	}

	return product, nil
}

func (ps *ProductService) CreateProduct(name, userID string) (*Product, *domain.Error) {
	id, err := util.GenerateID()
	if err != nil {
		return nil, err
	}

	product := &Product{
		ID:   id,
		Name: name,
		Owner: &users.User{
			ID: userID,
		},
	}

	return product, ps.ProductStorage.CreateProduct(product)
}
