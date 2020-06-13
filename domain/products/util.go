package products

func IsUserAuthorized(product *Product, userID string) bool {
	return product.Owner.ID == userID
}
