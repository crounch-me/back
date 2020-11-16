package lists

import (
	"time"

	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
)

// List represents a shopping list
type List struct {
	ID              string           `json:"id"`
	Name            string           `json:"name" validate:"required,lt=61"`
	CreationDate    time.Time        `json:"creationDate"`
	ArchivationDate *time.Time       `json:"archivationDate,omitempty"`
	Contributors    []*users.User    `json:"contributors,omitempty"`
	Products        []*ProductInList `json:"products,omitempty"`
}

// ProductInListLink represents a product in a list
type ProductInListLink struct {
	ProductID string `json:"productId"`
	ListID    string `json:"listId"`
	Buyed     bool   `json:"bought"`
}

type ProductInList struct {
	*products.Product
	Buyed bool `json:"bought"`
}

// UpdateProductInList represents the possible attributes to update in a product in a list
type UpdateProductInList struct {
	Buyed bool `json:"bought"`
}
