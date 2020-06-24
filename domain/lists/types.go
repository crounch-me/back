package lists

import (
	"time"

	"github.com/crounch-me/back/domain/users"
)

// List represents a shopping list
type List struct {
	ID           string      `json:"id"`
	Name         string      `json:"name" validate:"required,lt=61"`
	CreationDate time.Time   `json:"creationDate"`
	Owner        *users.User `json:"owner"`
}

// ProductInList represents a product in a list
type ProductInList struct {
	ProductID string `json:"productId"`
	ListID    string `json:"listId"`
}
