package domain

import "fmt"

type ErrorCallInfos struct {
	MethodName  string `json:"-"`
	PackageName string `json:"-"`
}

// Error describes the possible errors
type Error struct {
	Cause     error           `json:"-"`
	Code      string          `json:"error"`
	CallInfos *ErrorCallInfos `json:"-"`
}

const (
	ForbiddenErrorCode    = "forbidden-error"
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
func NewErrorWithCause(code string, cause error) *Error {
	return &Error{
		Code:  code,
		Cause: cause,
	}
}

func NewErrorWithCallInfosAndCause(code, packageName, methodName string, cause error) *Error {
	return &Error{
		Code:  code,
		Cause: cause,
		CallInfos: &ErrorCallInfos{
			MethodName:  methodName,
			PackageName: packageName,
		},
	}
}

// NewErrorWithCallInfos creates a new error with a code and the call informations
func NewErrorWithCallInfos(code, packageName, methodName string) *Error {
	return &Error{
		Code: code,
		CallInfos: &ErrorCallInfos{
			MethodName:  methodName,
			PackageName: packageName,
		},
	}
}

func (de *Error) Error() string {
	logString := fmt.Sprintf("'%s'", de.Code)

	if de.CallInfos != nil {
		logString += fmt.Sprintf(" appened in package '%s' and method '%s'", de.CallInfos.PackageName, de.CallInfos.MethodName)
	}

	if de.Cause != nil {
		logString += fmt.Sprintf(" because of '%s'", de.Cause)
	}

	return logString
}
