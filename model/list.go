package model

import "time"

// List represents a shopping list
type List struct {
	ID           string    `json:"id"`
	Name         string    `json:"name" validate:"required,lt=61"`
	CreationDate time.Time `json:"CreationDate"`
	Owner        *User
}

type ProductInList struct {
	ProductID string `json:"product_id"`
	ListID    string `json:"list_id"`
}
