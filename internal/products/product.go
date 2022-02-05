package products

import (
	"errors"

	"github.com/crounch-me/back/internal/common"
)

type Product struct {
	name string
}

func CreateProduct(name string) (Product, error) {
	if name == "" {
		return Product{}, errors.New(common.ERR_EMPTY_PRODUCT_NAME)
	}

	return Product{
		name: name,
	}, nil
}

func (p Product) Name() string {
	return p.name
}
