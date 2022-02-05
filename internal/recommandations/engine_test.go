package recommandations_test

import (
	"testing"

	"github.com/crounch-me/back/internal/baskets"
	"github.com/crounch-me/back/internal/recommandations"
	"github.com/stretchr/testify/assert"
)

func TestShampooRecommandation(t *testing.T) {
	// product_name := "Pomme de terre"
	// product, err := products.CreateProduct(product_name)
	// assert.Nil(t, err)

	basket_name := "Anniversaire de Raymond"
	basket_1, err := baskets.CreateBasket(basket_name)
	assert.Nil(t, err)

	baskets := []baskets.Basket{basket_1}
	next_product_recommandation_date := recommandations.Run(baskets)

	expected_recommandation_date := "Hello"
	assert.Equal(t, expected_recommandation_date, next_product_recommandation_date)
}
