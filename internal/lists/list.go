package lists

import (
	"errors"

	"github.com/crounch-me/back/internal/common"
	"github.com/crounch-me/back/internal/products"
)

type List struct {
	name     string
	articles []products.Product
}

func CreateList(name string) (List, error) {
	if name == "" {
		return List{}, errors.New(common.ERR_EMPTY_LIST_NAME)
	}
	return List{
		name: name,
	}, nil
}

func AddArticle(list List, product products.Product) (List, error) {
	new_articles := append(list.articles, product)

	return List{
		name:     list.name,
		articles: new_articles,
	}, nil
}

func CountArticles(list List) int {
	return len(list.articles)
}

func GetArticle(list List, index int) (products.Product, error) {
	if index < 0 || index >= len(list.articles) {
		return products.Product{}, errors.New(common.ERR_OUT_OF_RANGE_INDEX)
	}
	return list.articles[index], nil
}

func (l List) Name() string {
	return l.name
}
