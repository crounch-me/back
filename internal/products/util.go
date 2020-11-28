package products

import "github.com/crounch-me/back/domain/users"

func IsUserAuthorized(product *Product, userID string) bool {
	return product.Owner.ID == userID || product.Owner.ID == users.AdminID
}
