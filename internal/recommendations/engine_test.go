package recommendations_test

import (
	"testing"
	"time"

	"github.com/crounch-me/back/internal/baskets"
	"github.com/crounch-me/back/internal/common"
	"github.com/crounch-me/back/internal/recommendations"

	"github.com/stretchr/testify/assert"
)

const (
	article_id = "article_id"
	product_id = "product_id"
)

func TestIndexBoughtAtByProductIDOneArticleTwoDatesOK(t *testing.T) {
	now := time.Now()

	// create article
	article, err := baskets.CreateArticle(article_id, product_id)
	assert.Nil(t, err)

	// create first basket
	basket_1, err := baskets.CreateBasket("basket_name_1")
	assert.Nil(t, err)
	basket_1 = basket_1.AddArticle(article)
	assert.Nil(t, err)
	basket_1_finished_at := now.Add(-common.ONE_WEEK)
	basket_1 = basket_1.Finish(basket_1_finished_at)

	// create second basket
	basket_2, err := baskets.CreateBasket("basket_name_2")
	assert.Nil(t, err)
	basket_2 = basket_2.AddArticle(article)
	assert.Nil(t, err)
	basket_2_finished_at := now
	basket_2 = basket_2.Finish(basket_2_finished_at)

	// action
	all_baskets := []baskets.Basket{basket_1, basket_2}
	actual_bought_at := recommendations.IndexBoughtAtByArticle(all_baskets)

	// assert
	expected_bought_at := map[string][]time.Time{}
	expected_bought_at[article_id] = []time.Time{basket_1_finished_at, basket_2_finished_at}

	assert.Equal(t, expected_bought_at, actual_bought_at)
}

func TestComputeAverageBoughtDurationOneWeekBetweenTwoArticlesOK(t *testing.T) {
	now := time.Now()

	boughts_at := []time.Time{now, now.Add(common.ONE_WEEK)}

	average_bought_duration := recommendations.ComputeAverageBoughtDuration(boughts_at)

	assert.Equal(t, common.ONE_WEEK, average_bought_duration)
}

func TestComputeRecommendationDateByArticleOK(t *testing.T) {
	now := time.Now()

	// create article
	article, err := baskets.CreateArticle(article_id, product_id)
	assert.Nil(t, err)

	// create first basket
	basket_name_1 := "basket_name_1"
	basket_1, err := baskets.CreateBasket(basket_name_1)
	assert.Nil(t, err)
	basket_1 = basket_1.AddArticle(article)
	basket_1_finished_at := now.Add(-common.ONE_WEEK)
	basket_1 = basket_1.Finish(basket_1_finished_at)

	// create second basket
	basket_name_2 := "basket_name_2"
	basket_2, err := baskets.CreateBasket(basket_name_2)
	assert.Nil(t, err)
	basket_2 = basket_2.AddArticle(article)
	basket_2_finished_at := now
	basket_2 = basket_2.Finish(basket_2_finished_at)

	// action
	all_baskets := []baskets.Basket{basket_1, basket_2}
	articles_to_buy := recommendations.ComputeRecommendationDateByArticle(all_baskets)

	// assert
	expected_articles := make(map[string]time.Time, 0)
	expected_articles[article_id] = now.Add(common.ONE_WEEK)

	assert.Equal(t, expected_articles, articles_to_buy)
}
