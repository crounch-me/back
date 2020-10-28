package lists

import (
	"time"

	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
)

// List represents a shopping list
type List struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name" validate:"required,lt=61"`
	CreationDate time.Time                `json:"creationDate"`
	Owner        *users.User              `json:"owner,omitempty"`
	Products     []*ProductInList `json:"products,omitempty"`
}

// ProductInListLink represents a product in a list
type ProductInListLink struct {
	ProductID string `json:"productId"`
	ListID    string `json:"listId"`
	Buyed     bool   `json:"buyed"`
}

type ProductInList struct {
	*products.Product
	Buyed bool `json:"buyed"`
}

// UpdateProductInList represents the possible attributes to update in a product in a list
type UpdateProductInList struct {
	Buyed bool `json:"buyed"`
}
