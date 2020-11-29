package lists

import "errors"

type Product struct {
	uuid string
}

var (
	ErrProductAlreadyInList = errors.New("product already in list")
)

func NewProduct(uuid string) (*Product, error) {
	if uuid == "" {
		return nil, errors.New("empty product uuid")
	}

	return &Product{
		uuid: uuid,
	}, nil
}

func (p Product) UUID() string {
	return p.uuid
}

func (l *List) AddProduct(p *Product) error {
	if l.HasProduct(p.uuid) {
		return ErrProductAlreadyInList
	}

	l.products = append(l.products, p)

	return nil
}
