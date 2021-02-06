package app

import (
	"errors"

	"github.com/crounch-me/back/internal/common/utils"
	"github.com/crounch-me/back/internal/products/domain/products"
)

type ProductService struct {
	generationLibrary  utils.GenerationLibrary
	productsRepository Repository
}

func NewProductService(generationLibrary utils.GenerationLibrary, productsRepository Repository) (*ProductService, error) {
	if generationLibrary == nil {
		return nil, errors.New("generationLibrary is nil")
	}

	if productsRepository == nil {
		return nil, errors.New("productsRepository is nil")
	}

	return &ProductService{
		generationLibrary:  generationLibrary,
		productsRepository: productsRepository,
	}, nil
}

func (s *ProductService) CreateProduct(name string) (string, error) {
	productUUID, err := s.generationLibrary.UUID()
	if err != nil {
		return "", err
	}

	p, err := products.NewProduct(productUUID, name, nil)
	if err != nil {
		return "", err
	}

	err = s.productsRepository.SaveProduct(p)
	if err != nil {
		return "", err
	}

	return productUUID, nil
}

func (s *ProductService) SearchDefaults(name string) ([]*products.Product, error) {
	return s.productsRepository.SearchDefaults(name)
}
