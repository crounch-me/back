package lists

import "errors"

var (
	ErrUserNotContributor = errors.New("user is not a contributor of the list")

	ErrListNotFound = errors.New("list not found")

	ErrProductAlreadyInList  = errors.New("product already in list")
	ErrProductNotFoundInList = errors.New("product not found in list")
)
