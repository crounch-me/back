package baskets_test

import (
	"testing"
	"time"

	"github.com/crounch-me/back/internal/baskets"
	"github.com/crounch-me/back/internal/common"
	"github.com/stretchr/testify/assert"
)

const (
	valid_article_id_1 = "article id 1"
	valid_article_id_2 = "article id 2"
	valid_basket_name  = "Anniversaire de Raymond"
	valid_product_id_1 = "product id 1"
	valid_product_id_2 = "product id 2"
)

var test_article_1, _ = baskets.CreateArticle(valid_article_id_1, valid_product_id_1)
var test_article_2, _ = baskets.CreateArticle(valid_article_id_2, valid_product_id_2)
var test_basket, _ = baskets.CreateBasket(valid_basket_name)

func TestCreateBasketEmptyName(t *testing.T) {
	name := ""
	_, err := baskets.CreateBasket(name)

	assert.Equal(t, common.ERR_EMPTY_BASKET_NAME, err.Error())
}

func TestCreateBasketOK(t *testing.T) {
	basket, err := baskets.CreateBasket(valid_basket_name)

	assert.Nil(t, err)
	assert.Equal(t, valid_basket_name, basket.Name())
}

func TestAddArticle(t *testing.T) {
	article, err := baskets.CreateArticle(valid_article_id_1, valid_product_id_1)
	assert.Nil(t, err)

	basket, err := baskets.CreateBasket(valid_basket_name)
	assert.Nil(t, err)

	new_basket, err := basket.AddArticle(article)
	assert.Nil(t, err)
	assert.Equal(t, 1, new_basket.CountArticles())
}

func TestFinishOK(t *testing.T) {
	finish_time := time.Now()

	new_basket := test_basket.Finish(finish_time)

	assert.Equal(t, finish_time, new_basket.FinishedAt())
}
