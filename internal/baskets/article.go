package baskets

import (
	"errors"

	"github.com/crounch-me/back/internal/common"
)

type Article struct {
	id         string
	product_id string
}

func CreateArticle(id, product_id string) (Article, error) {
	if id == "" {
		return Article{}, errors.New(common.ERR_EMPTY_ARTICLE_ID)
	}

	if product_id == "" {
		return Article{}, errors.New(common.ERR_EMPTY_PRODUCT_ID)
	}

	return Article{
		id:         id,
		product_id: product_id,
	}, nil
}

func (a Article) ID() string {
	return a.id
}

func (a Article) ProductID() string {
	return a.product_id
}

func (b Basket) AddArticle(article Article) (Basket, error) {
	new_articles := append(b.articles, article)

	b.articles = new_articles
	return b, nil
}

func (b Basket) ForEachArticle(callback func(article Article) error) error {
	for _, article := range b.articles {
		err := callback(article)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b Basket) CountArticles() int {
	return len(b.articles)
}
