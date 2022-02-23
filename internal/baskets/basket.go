package baskets

import (
	"errors"
	"time"

	"github.com/crounch-me/back/internal/common"
)

type Basket struct {
	name        string
	finished_at time.Time
	articles    []Article
}

func CreateBasket(name string) (Basket, error) {
	if name == "" {
		return Basket{}, errors.New(common.ERR_EMPTY_BASKET_NAME)
	}

	return Basket{
		name: name,
	}, nil
}

func (l Basket) Name() string {
	return l.name
}

func (l Basket) FinishedAt() time.Time {
	return l.finished_at
}

func (b Basket) Finish(finished_at time.Time) Basket {
	b.finished_at = finished_at
	return b
}
