package products_test

import (
	"testing"

	"github.com/crounch-me/back/internal/common"
	"github.com/crounch-me/back/internal/products"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductEmptyName(t *testing.T) {
	name := ""
	_, err := products.CreateProduct(name)

	assert.Equal(t, common.ERR_EMPTY_PRODUCT_NAME, err.Error())
}

func TestCreateProductOK(t *testing.T) {
	name := "Pomme de terre"
	product, err := products.CreateProduct(name)

	assert.Equal(t, name, product.Name())
	assert.Nil(t, err)
}
