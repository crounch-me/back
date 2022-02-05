package baskets

import (
	"errors"

	"github.com/crounch-me/back/internal/common"
	"github.com/crounch-me/back/internal/products"
)

type Basket struct {
	name     string
	articles []products.Product
}

func CreateBasket(name string) (Basket, error) {
	if name == "" {
		return Basket{}, errors.New(common.ERR_EMPTY_BASKET_NAME)
	}
	return Basket{
		name: name,
	}, nil
}

func AddArticle(basket Basket, product products.Product) (Basket, error) {
	new_articles := append(basket.articles, product)

	return Basket{
		name:     basket.name,
		articles: new_articles,
	}, nil
}

func CountArticles(basket Basket) int {
	return len(basket.articles)
}

func GetArticle(basket Basket, index int) (products.Product, error) {
	if index < 0 || index >= len(basket.articles) {
		return products.Product{}, errors.New(common.ERR_OUT_OF_RANGE_INDEX)
	}
	return basket.articles[index], nil
}

func (l Basket) Name() string {
	return l.name
}
