package baskets_test

import (
	"testing"

	"github.com/crounch-me/back/internal/baskets"
	"github.com/crounch-me/back/internal/common"
	"github.com/stretchr/testify/assert"
)

const (
	article_id = "article_id"
	product_id = "product_id"
)

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
