package lists

import "errors"

type Product struct {
	uuid string
}

func NewProduct(uuid string) (*Product, error) {
	if uuid == "" {
		return nil, errors.New("empty product uuid")
	}

	return &Product{
		uuid: uuid,
	}, nil
}

func (l *List) AddProduct(p *Product) error {
	if l.HasProduct(p.uuid) {
		return ErrProductAlreadyInList
	}

	l.products = append(l.products, p)

	return nil
}
