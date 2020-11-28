package internal

import "fmt"

type CallError struct {
	MethodName  string `json:"-"`
	PackageName string `json:"-"`
}

type FieldError struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

// Error describes the possible errors
type Error struct {
	Cause  error         `json:"-"`
	Code   string        `json:"error"`
	Call   *CallError    `json:"-"`
	Fields []*FieldError `json:"fields,omitempty"`
}

const (
	ForbiddenErrorCode    = "forbidden-error"
	UnauthorizedErrorCode = "unauthorized-error"
	UnknownErrorCode      = "unknown-error"
	UnmarshalErrorCode    = "unmarshal-error"
	InvalidErrorCode      = "invalid-error"
)

// NewError creates a new error from the /internal
func NewError(code string) *Error {
	return &Error{
		Code: code,
	}
}

func (e *Error) WithCause(cause error) *Error {
	e.Cause = cause
	return e
}

func (e *Error) WithCallError(callError *CallError) *Error {
	e.Call = callError
	return e
}

func (e *Error) WithCall(packageName, methodName string) *Error {
	e.Call = &CallError{
		MethodName:  methodName,
		PackageName: packageName,
	}
	return e
}

func (e *Error) WithFields(fields []*FieldError) *Error {
	e.Fields = fields
	return e
}

// NewErrorWithCallInfos creates a new error with a code and the call informations
func NewErrorWithCallInfos(code, packageName, methodName string) *Error {
	return &Error{
		Code: code,
		Call: &CallError{
			MethodName:  methodName,
			PackageName: packageName,
		},
	}
}

func (de *Error) Error() string {
	logString := fmt.Sprintf("'%s'", de.Code)

	if de.Call != nil {
		logString += fmt.Sprintf(" appened in package '%s' and method '%s'", de.Call.PackageName, de.Call.MethodName)
	}

	if de.Cause != nil {
		logString += fmt.Sprintf(" because of '%s'", de.Cause)
	}

	return logString
}
