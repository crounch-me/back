package lists_test

import (
	"testing"

	"github.com/crounch-me/back/internal/common"
	"github.com/crounch-me/back/internal/lists"
	"github.com/crounch-me/back/internal/products"
	"github.com/stretchr/testify/assert"
)

const (
	valid_product_name_1 = "Pomme de terre"
	valid_product_name_2 = "Courgette"
	valid_list_name      = "Anniversaire de Raymond"
)

var test_product_1, _ = products.CreateProduct(valid_product_name_1)
var test_product_2, _ = products.CreateProduct(valid_product_name_2)
var test_list, _ = lists.CreateList(valid_list_name)

func TestCreateListEmptyName(t *testing.T) {
	name := ""
	_, err := lists.CreateList(name)

	assert.Equal(t, common.ERR_EMPTY_LIST_NAME, err.Error())
}

func TestCreateListOK(t *testing.T) {
	name := "Anniversaire de Raymond"
	list, err := lists.CreateList(name)

	assert.Equal(t, name, list.Name())
	assert.Nil(t, err)
}

func TestAddProduct(t *testing.T) {
	product_name := "Pomme de terre"
	product, err := products.CreateProduct(product_name)
	assert.Nil(t, err)

	list_name := "Anniversaire de Raymond"
	list, err := lists.CreateList(list_name)
	assert.Nil(t, err)

	new_list, err := lists.AddArticle(list, product)
	assert.Equal(t, 1, lists.CountArticles(new_list))
	// assert.Equal(t, new_list.Get(0).name, product.name)
	assert.Nil(t, err)
}

func TestGetArticleNegativeIndex(t *testing.T) {
	_, err := lists.GetArticle(test_list, -1)

	assert.Equal(t, common.ERR_OUT_OF_RANGE_INDEX, err.Error())
}

func TestGetArticleTooHighIndex(t *testing.T) {
	_, err := lists.GetArticle(test_list, 1)

	assert.Equal(t, common.ERR_OUT_OF_RANGE_INDEX, err.Error())
}

func TestGetArticleOK(t *testing.T) {
	list_with_one_article, err := lists.AddArticle(test_list, test_product_1)
	assert.Nil(t, err)
	list_with_two_articles, err := lists.AddArticle(list_with_one_article, test_product_2)
	assert.Nil(t, err)

	article_1, err := lists.GetArticle(list_with_two_articles, 0)
	assert.Nil(t, err)
	assert.Equal(t, valid_product_name_1, article_1.Name())

	article_2, err := lists.GetArticle(list_with_two_articles, 1)
	assert.Nil(t, err)
	assert.Equal(t, valid_product_name_2, article_2.Name())
}
