package adapters

import (
	"errors"
	"strings"
	"sync"
	"unicode"

	"github.com/crounch-me/back/internal/common/utils"
	"github.com/crounch-me/back/internal/products/app"
	"github.com/crounch-me/back/internal/products/domain/products"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type ProductsMemoryRepository struct {
	products   map[string]*products.Product
	categories map[string]*products.Category
	lock       *sync.RWMutex
	generation utils.GenerationLibrary
}

func NewProductsMemoryRepository(generation utils.GenerationLibrary) (app.Repository, error) {
	if generation == nil {
		return nil, errors.New("generation is nil")
	}

	return &ProductsMemoryRepository{
		products:   make(map[string]*products.Product, 0),
		categories: make(map[string]*products.Category, 0),
		generation: generation,
		lock:       &sync.RWMutex{},
	}, nil
}

func (r *ProductsMemoryRepository) InsertDefaultValues() error {
	epicerie, err := r.createCategory("Epicerie")
	if err != nil {
		return err
	}

	boucherie, err := r.createCategory("Boucherie")
	if err != nil {
		return err
	}

	_, err = r.createProduct("Lentille", epicerie)
	if err != nil {
		return err
	}

	_, err = r.createProduct("Saucisse de Morteau", boucherie)
	if err != nil {
		return err
	}

	_, err = r.createProduct("Saucisse de Montbéliard", boucherie)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductsMemoryRepository) SaveProduct(p *products.Product) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.products[p.UUID()]; ok {
		return errors.New("product already exists")
	}

	if p.HasCategory() {
		if _, ok := r.categories[p.CategoryUUID()]; !ok {
			return errors.New("category not found")
		}
	}

	r.products[p.UUID()] = p

	return nil
}

func (r *ProductsMemoryRepository) ReadByID(uuid string) (*products.Product, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if p, ok := r.products[uuid]; ok {
		return p, nil
	}

	return nil, products.ErrProductNotFound
}

func (r *ProductsMemoryRepository) SearchDefaults(name string) ([]*products.Product, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	results := make([]*products.Product, 0)
	formattedNameToSearch, err := removeAccentuatedCharaceters(strings.ToLower(name))
	if err != nil {
		return nil, err
	}

	for _, product := range r.products {
		formattedNameOfProduct, err := removeAccentuatedCharaceters(strings.ToLower(product.Name()))
		if err != nil {
			return nil, err
		}

		if strings.Contains(formattedNameOfProduct, formattedNameToSearch) {
			results = append(results, product)
		}
	}

	return results, nil
}

func (r *ProductsMemoryRepository) saveCategory(c *products.Category) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.categories[c.UUID()]; ok {
		return errors.New("category already exists")
	}

	r.categories[c.UUID()] = c

	return nil
}

func (r *ProductsMemoryRepository) createProduct(name string, category *products.Category) (*products.Product, error) {
	id, err := r.generation.UUID()
	if err != nil {
		return nil, err
	}

	product, err := products.NewProduct(id, name, category)
	if err != nil {
		return nil, err
	}

	err = r.SaveProduct(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductsMemoryRepository) createCategory(name string) (*products.Category, error) {
	id, err := r.generation.UUID()
	if err != nil {
		return nil, err
	}

	category, err := products.NewCategory(id, name)
	if err != nil {
		return nil, err
	}

	err = r.saveCategory(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func removeAccentuatedCharaceters(str string) (string, error) {
	transf := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

	output, _, err := transform.String(transf, str)
	if err != nil {
		return "", err
	}

	return output, nil
}
