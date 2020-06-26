package products

import (
	"github.com/crounch-me/back/domain/categories"
	"github.com/crounch-me/back/domain/users"
)

type Product struct {
	ID       string               `json:"id"`
	Name     string               `json:"name" validate:"required,lt=61"`
	Owner    *users.User          `json:"owner,omitempty"`
	Category *categories.Category `json:"category"`
}
