package lists

import "errors"

type Product struct {
	uuid   string
	bought bool
}

var (
	ErrProductAlreadyInList = errors.New("product already in list")
	ErrProductNotFound      = errors.New("product not found in list")

	NotFoundIndex = -1
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

func (l List) HasProduct(uuid string) bool {
	for _, p := range l.products {
		if p.uuid == uuid {
			return true
		}
	}

	return false
}

func (l *List) AddProduct(p *Product) error {
	if l.HasProduct(p.uuid) {
		return ErrProductAlreadyInList
	}

	l.products = append(l.products, p)

	return nil
}

func (l List) Products() []Product {
	products := make([]Product, 0)

	for _, product := range l.products {
		products = append(products, *product)
	}

	return products
}

func (l List) FindProductIndex(uuid string) (int, error) {
	productIndex := NotFoundIndex
	for i, p := range l.products {
		if p.UUID() == uuid {
			productIndex = i
		}
	}

	if productIndex == NotFoundIndex {
		return 0, ErrProductNotFound
	}

	return productIndex, nil
}

func (l *List) RemoveProduct(uuid string) error {
	productIndex, err := l.FindProductIndex(uuid)
	if err != nil {
		return err
	}

	l.products = remove(l.products, productIndex)

	return nil
}

func (p *Product) Buy() {
	p.bought = true
}

func remove(products []*Product, i int) []*Product {
	products[i] = products[len(products)-1]

	return products[:len(products)-1]
}
