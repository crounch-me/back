package products

import (
	"github.com/crounch-me/back/internal/account"
	"github.com/crounch-me/back/internal/categories"
)

type Product struct {
	ID       string               `json:"id"`
	Name     string               `json:"name" validate:"required,lt=61"`
	Owner    *account.User        `json:"owner,omitempty"`
	Category *categories.Category `json:"category,omitempty"`
}
