package products

import "github.com/crounch-me/back/internal/user"

func IsUserAuthorized(product *Product, userID string) bool {
	return product.Owner.ID == userID || product.Owner.ID == user.AdminID
}
