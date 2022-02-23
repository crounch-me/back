package baskets_test

import (
	"errors"
	"testing"

	"github.com/crounch-me/back/internal/baskets"
	"github.com/crounch-me/back/internal/common"
	"github.com/stretchr/testify/assert"
)

const (
	article_id = "article_id"
	product_id = "product_id"
)

func TestForEachArticleOK(t *testing.T) {
	basket, err := baskets.CreateBasket("basket name")
	assert.Nil(t, err)

	article, err := baskets.CreateArticle(article_id, product_id)
	assert.Nil(t, err)

	basket = basket.AddArticle(article)

	i := 0
	err = basket.ForEachArticle(func(article baskets.Article) error {
		i++
		return nil
	})

	assert.Nil(t, err)
	assert.Equal(t, basket.CountArticles(), i)
}

func TestForEachArticleError(t *testing.T) {
	basket, err := baskets.CreateBasket("basket name")
	assert.Nil(t, err)

	article, err := baskets.CreateArticle(article_id, product_id)
	assert.Nil(t, err)

	error_message := "unknown-error"

	basket = basket.AddArticle(article)

	returned_err := basket.ForEachArticle(func(article baskets.Article) error {
		return errors.New(error_message)
	})

	assert.Equal(t, error_message, returned_err.Error())
}

func TestCreateArticleEmptyID(t *testing.T) {
	_, err := baskets.CreateArticle("", product_id)

	assert.Equal(t, common.ERR_EMPTY_ARTICLE_ID, err.Error())
}

func TestCreateArticleEmptyProductID(t *testing.T) {
	_, err := baskets.CreateArticle(article_id, "")

	assert.Equal(t, common.ERR_EMPTY_PRODUCT_ID, err.Error())
}

func TestCreateArticleOK(t *testing.T) {
	article, err := baskets.CreateArticle(article_id, product_id)

	assert.Nil(t, err)
	assert.Equal(t, article_id, article.ID())
	assert.Equal(t, product_id, article.ProductID())
}
