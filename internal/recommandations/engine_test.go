package recommandations_test

import (
	"testing"

	"github.com/crounch-me/back/internal/lists"
	"github.com/crounch-me/back/internal/recommandations"
	"github.com/stretchr/testify/assert"
)

func TestShampooRecommandation(t *testing.T) {
	// product_name := "Pomme de terre"
	// product, err := products.CreateProduct(product_name)
	// assert.Nil(t, err)

	list_name := "Anniversaire de Raymond"
	list_1, err := lists.CreateList(list_name)
	assert.Nil(t, err)

	lists := []lists.List{list_1}
	next_product_recommandation_date := recommandations.Run(lists)

	expected_recommandation_date := "Hello"
	assert.Equal(t, expected_recommandation_date, next_product_recommandation_date)
}
