package adapters

import (
	"errors"
	"sync"

	"github.com/crounch-me/back/internal/products/app"
	"github.com/crounch-me/back/internal/products/domain/products"
)

type ProductsMemoryRepository struct {
	products map[string]*products.Product
	lock     *sync.RWMutex
}

func NewProductsMemoryRepository() app.Repository {
	return &ProductsMemoryRepository{
		products: make(map[string]*products.Product, 0),
		lock:     &sync.RWMutex{},
	}
}

func (r *ProductsMemoryRepository) SaveProduct(p *products.Product) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.products[p.UUID()]; ok {
		return errors.New("product already exists")
	}

	r.products[p.UUID()] = p

	return nil
}

func (r *ProductsMemoryRepository) SearchDefaults(name string) ([]*products.Product, error) {
	return []*products.Product{}, nil
}
