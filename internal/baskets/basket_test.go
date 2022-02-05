package baskets_test

import (
	"testing"
	"time"

	"github.com/crounch-me/back/internal/baskets"
	"github.com/crounch-me/back/internal/common"
	"github.com/crounch-me/back/internal/products"
	"github.com/stretchr/testify/assert"
)

const (
	valid_product_name_1 = "Pomme de terre"
	valid_product_name_2 = "Courgette"
	valid_basket_name    = "Anniversaire de Raymond"
)

var test_product_1, _ = products.CreateProduct(valid_product_name_1)
var test_product_2, _ = products.CreateProduct(valid_product_name_2)
var test_basket, _ = baskets.CreateBasket(valid_basket_name)

func TestCreateBasketEmptyName(t *testing.T) {
	name := ""
	_, err := baskets.CreateBasket(name)

	assert.Equal(t, common.ERR_EMPTY_BASKET_NAME, err.Error())
}

func TestCreateBasketOK(t *testing.T) {
	name := "Anniversaire de Raymond"
	basket, err := baskets.CreateBasket(name)

	assert.Equal(t, name, basket.Name())
	assert.Nil(t, err)
}

func TestAddArticle(t *testing.T) {
	product, err := products.CreateProduct(valid_product_name_1)
	assert.Nil(t, err)

	basket, err := baskets.CreateBasket(valid_basket_name)
	assert.Nil(t, err)

	new_basket, err := baskets.AddArticle(basket, product)
	assert.Equal(t, 1, baskets.CountArticles(new_basket))
	// assert.Equal(t, new_basket.Get(0).name, product.name)
	assert.Nil(t, err)
}

func TestGetArticleNegativeIndex(t *testing.T) {
	_, err := baskets.GetArticle(test_basket, -1)

	assert.Equal(t, common.ERR_OUT_OF_RANGE_INDEX, err.Error())
}

func TestGetArticleTooHighIndex(t *testing.T) {
	_, err := baskets.GetArticle(test_basket, 1)

	assert.Equal(t, common.ERR_OUT_OF_RANGE_INDEX, err.Error())
}

func TestGetArticleOK(t *testing.T) {
	basket_with_one_article, err := baskets.AddArticle(test_basket, test_product_1)
	assert.Nil(t, err)
	basket_with_two_articles, err := baskets.AddArticle(basket_with_one_article, test_product_2)
	assert.Nil(t, err)

	article_1, err := baskets.GetArticle(basket_with_two_articles, 0)
	assert.Nil(t, err)
	assert.Equal(t, valid_product_name_1, article_1.Name())

	article_2, err := baskets.GetArticle(basket_with_two_articles, 1)
	assert.Nil(t, err)
	assert.Equal(t, valid_product_name_2, article_2.Name())
}

func TestFinishOK(t *testing.T) {
	finish_time := time.Now()
	new_basket, err := baskets.Finish(test_basket, finish_time)
	assert.Nil(t, err)
	assert.Equal(t, finish_time, new_basket.IsFinishedAt())
}
