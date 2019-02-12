package model

import "fmt"

type Type int

const (
	ErrNotFound Type = iota
	ErrCreation
	ErrWrongPassword
)

type DatabaseError struct {
	Cause error
	Type  Type
}

func NewDatabaseError(t Type, cause error) error {
	return &DatabaseError{
		Type:  t,
		Cause: cause,
	}
}

func (de *DatabaseError) Error() string {
	if de.Cause != nil {
		return fmt.Sprintf("type %d, cause %s", de.Type, de.Cause.Error())
	}
	return fmt.Sprintf("type %d, no cause", de.Type)
}
