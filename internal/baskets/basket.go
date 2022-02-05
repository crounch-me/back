package baskets

import (
	"errors"
	"time"

	"github.com/crounch-me/back/internal/common"
	"github.com/crounch-me/back/internal/products"
)

type Basket struct {
	name        string
	finished_at time.Time
	articles    []products.Product
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

	new_basket := basket
	new_basket.articles = new_articles
	return new_basket, nil
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

func (l Basket) IsFinishedAt() time.Time {
	return l.finished_at
}

func Finish(basket Basket, finished_at time.Time) (Basket, error) {
	new_basket := basket
	new_basket.finished_at = finished_at
	return new_basket, nil
}
