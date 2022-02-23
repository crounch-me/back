package recommandations

import (
	"log"
	"time"

	"github.com/crounch-me/back/internal/baskets"
)

func IndexBoughtAtByProductID(all_baskets []baskets.Basket) map[string][]time.Time {
	bought_at_indexed_by_product_id := make(map[string][]time.Time)

	for _, basket := range all_baskets {
		log.Println("basket")
		basket.ForEachArticle(func(article baskets.Article) error {
			log.Println("article")
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

func RecommandArticles(all_baskets []baskets.Basket) map[string]time.Time {
	return map[string]time.Time{}
}
