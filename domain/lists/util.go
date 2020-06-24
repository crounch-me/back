package lists

func IsUserAuthorized(list *List, userID string) bool {
	return list.Owner.ID == userID
}
