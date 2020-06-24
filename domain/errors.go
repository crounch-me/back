package domain

import "fmt"

// Error describes the possible errors
type Error struct {
	Code  string `json:"error"`
	Cause error  `json:"-"`
}

const (
	UnauthorizedErrorCode = "unauthorized-error"
	UnknownErrorCode      = "unknown-error"
	UnmarshalErrorCode    = "unmarshal-error"
	InvalidErrorCode      = "invalid-error"
)

// NewError creates a new error from the domain
func NewError(code string) *Error {
	return &Error{
		Code: code,
	}
}

// NewErrorWithCause creates a new error from the domain with a cause error
func NewErrorWithCause(code string, err error) *Error {
	return &Error{
		Code:  code,
		Cause: err,
	}
}

func (de *Error) Error() string {
	if de.Cause == nil {
		return fmt.Sprintf("%s", de.Code)
	}

	return fmt.Sprintf("'%s' because of '%s'", de.Code, de.Cause)
}
