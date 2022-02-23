package products_test

import (
	"testing"

	"github.com/crounch-me/back/internal/common"
	"github.com/crounch-me/back/internal/products"
	"github.com/stretchr/testify/assert"
)

const (
	valid_id   = "valid id"
	valid_name = "Pomme de terre"
)

func TestCreateProductEmptyId(t *testing.T) {
	id := ""
	_, err := products.CreateProduct(id, "unknown")
	assert.Equal(t, common.ERR_EMPTY_PRODUCT_ID, err.Error())
}

func TestCreateProductEmptyName(t *testing.T) {
	name := ""
	_, err := products.CreateProduct(valid_id, name)

	assert.Equal(t, common.ERR_EMPTY_PRODUCT_NAME, err.Error())
}

func TestCreateProductOK(t *testing.T) {
	product, err := products.CreateProduct(valid_id, valid_name)

	assert.Equal(t, valid_name, product.Name())
	assert.Equal(t, valid_id, product.ID())
	assert.Nil(t, err)
}
