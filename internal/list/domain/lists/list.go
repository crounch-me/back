package lists

import (
	"errors"
	"time"
)

type List struct {
	uuid            string
	name            string
	creationDate    time.Time
	archivationDate *time.Time
	contributors    []string
	products        []*Product
}

func NewList(uuid string, name string, creationDate time.Time, archivationDate *time.Time) (*List, error) {
	if uuid == "" {
		return nil, errors.New("empty list uuid")
	}

	if name == "" {
		return nil, errors.New("empty list name")
	}

	if creationDate.IsZero() {
		return nil, errors.New("empty list creation date")
	}

	return &List{
		uuid:            uuid,
		name:            name,
		creationDate:    creationDate,
		archivationDate: archivationDate,
		products:        []*Product{},
		contributors:    []string{},
	}, nil
}

func (l List) UUID() string {
	return l.uuid
}

func (l List) Name() string {
	return l.name
}

func (l List) CreationDate() time.Time {
	return l.creationDate
}

func (l List) Contributors() []string {
	return l.contributors
}

func (l List) Products() []Product {
	products := make([]Product, 0)

	for _, product := range l.products {
		products = append(products, *product)
	}

	return products
}

func (l List) HasProduct(uuid string) bool {
	for _, p := range l.products {
		if p.uuid == uuid {
			return true
		}
	}

	return false
}
