package recommendations

import (
	"time"

	"github.com/crounch-me/back/internal/baskets"
	"github.com/crounch-me/back/internal/common"
)

func IndexBoughtAtByArticle(all_baskets []baskets.Basket) map[string][]time.Time {
	bought_at_indexed_by_product_id := make(map[string][]time.Time)

	for _, basket := range all_baskets {
		basket.ForEachArticle(func(article baskets.Article) error {
			boughts_at, ok := bought_at_indexed_by_product_id[article.ID()]
			if ok {
				bought_at_indexed_by_product_id[article.ID()] = append(boughts_at, basket.FinishedAt())
			} else {
				bought_at_indexed_by_product_id[article.ID()] = []time.Time{basket.FinishedAt()}
			}

			return nil
		})
	}

	return bought_at_indexed_by_product_id
}

func ComputeAverageBoughtDuration(boughts_at []time.Time) time.Duration {
	durations_sum := time.Duration(0)
	duration_count := time.Duration(len(boughts_at) - 1)
	for i := 1; i < len(boughts_at); i++ {
		durations_sum += boughts_at[i].Sub(boughts_at[i-1])
	}

	return time.Duration(durations_sum / duration_count)
}

func ComputeRecommendationDateByArticle(all_baskets []baskets.Basket) map[string]time.Time {
	boughts_at_indexed_by_article_id := IndexBoughtAtByArticle(all_baskets)
	recommendation_dates := make(map[string]time.Time, len(boughts_at_indexed_by_article_id))

	for article_id, boughts_at := range boughts_at_indexed_by_article_id {
		average_duration := ComputeAverageBoughtDuration(boughts_at)
		recommendation_dates[article_id] = boughts_at[len(boughts_at)-1].Add(average_duration)
	}

	return recommendation_dates
}

func FilterRecommendedArticles(recommended_articles map[string]time.Time, start, end time.Time) []string {
	filtered_articles := []string{}

	for article_id, recommendation_date := range recommended_articles {
		if common.IsInRange(recommendation_date, start, end) {
			filtered_articles = append(filtered_articles, article_id)
		}
	}

	return filtered_articles
}
