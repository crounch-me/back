package products

import "errors"

type Product struct {
	uuid     string
	name     string
	category *Category
}

func NewProduct(uuid, name string, category *Category) (*Product, error) {
	if uuid == "" {
		return nil, errors.New("product uuid is empty")
	}

	if name == "" {
		return nil, errors.New("product name is empty")
	}

	if category == nil {
		return nil, errors.New("product category is empty")
	}

	return &Product{
		uuid:     uuid,
		name:     name,
		category: category,
	}, nil
}

func (p *Product) UUID() string {
	return p.uuid
}
