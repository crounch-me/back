package errorcode

const (
	DuplicateCode        = "duplicate"
	DuplicateDescription = "Entity already exists"

	NotFoundCode = "not-found"

	WrongPasswordCode        = "wrong-password"
	WrongPasswordDescription = "The password doesn't match the saved one"

	InvalidCode        = "invalid"
	InvalidDescription = "The field %s has failed validation, reason: '%s', actualValue: '%v'"

	DatabaseCode        = "database"
	DatabaseDescription = "An error occured while accessing the database"

	UnmarshalCode        = "unmarshal"
	UnmarshalDescription = "An error occured while reading request body"

	UserDataCode        = "user-data"
	UserDataDescription = "An error occured while retrieving user data"
)
