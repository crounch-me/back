package products

import "github.com/crounch-me/back/internal/account"

func IsUserAuthorized(product *Product, userID string) bool {
	return product.Owner.ID == userID || product.Owner.ID == account.AdminID
}
