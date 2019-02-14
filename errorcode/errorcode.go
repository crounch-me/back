package errorcode

const (
	Duplicate            = "duplicate"
	DuplicateDescription = "User with this email already exists"

	NotFound = "not-found"

	WrongPasswordCode        = "wrong-password"
	WrongPasswordDescription = "The password doesn't match the saved one"

	InvalidCode        = "invalid"
	InvalidDescription = "The field %s has failed validation, reason: '%s'"

	DatabaseCode        = "database"
	DatabaseDescription = "An error occured while accessing the database"

	UnmarshalCode        = "unmarshal"
	UnmarshalDescription = "An error occured while reading request body"
)
