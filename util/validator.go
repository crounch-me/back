package util

import (
	"reflect"
	"strings"

	"github.com/crounch-me/back/domain"
	"gopkg.in/go-playground/validator.v9"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	v := &Validator{
		validate: validate,
	}

	return v
}

func (v *Validator) Struct(s interface{}) error {
	return v.validate.Struct(s)
}

func (v *Validator) Var(fieldName, value, tag string) *domain.Error {
	err := v.validate.Var(value, tag)
	if err != nil {
		fields := []*domain.FieldError{
			{
				Name:  fieldName,
				Error: tag,
			},
		}
		return domain.NewError(domain.InvalidErrorCode).WithFields(fields)
	}

	return nil
}
