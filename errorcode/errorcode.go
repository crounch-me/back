package errorcode

const (
	Duplicate = "duplicate"
	NotFound  = "not-found"

	DatabaseCode        = "database"
	DatabaseDescription = "An error occured while accessing the database"

	UnmarshalCode        = "unmarshal"
	UnmarshalDescription = "An error occured while reading request body"
)
