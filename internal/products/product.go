package products

import (
	"errors"

	"github.com/crounch-me/back/internal/common"
)

type Product struct {
	id   string
	name string
}

func CreateProduct(id string, name string) (Product, error) {
	if id == "" {
		return Product{}, errors.New(common.ERR_EMPTY_PRODUCT_ID)
	}

	if name == "" {
		return Product{}, errors.New(common.ERR_EMPTY_PRODUCT_NAME)
	}

	return Product{
		id:   id,
		name: name,
	}, nil
}

func (p Product) Name() string {
	return p.name
}

func (p Product) ID() string {
	return p.id
}
