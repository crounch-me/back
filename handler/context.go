package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/crounch-me/back/configuration"
	"github.com/crounch-me/back/errorcode"
	"github.com/crounch-me/back/model"
	"github.com/crounch-me/back/storage"
	"github.com/crounch-me/back/storage/mock"
	"github.com/crounch-me/back/storage/postgres"
)

const (
	// ContextUserID is the key to retrieve user id from gin context
	ContextUserID = "ContextUserID"
)

// Context holds everything to respond to requests
type Context struct {
	Storage  storage.Storage
	Config   *configuration.Config
	Validate *validator.Validate
}

// NewContext creates and initialize everything for the requests
func NewContext(config *configuration.Config) *Context {
	var storage storage.Storage

	if config.Mock {
		storage = mock.NewStorageMock()
	} else {
		postgres.InitDB(config.DBConnectionURI)
		storage = postgres.NewStorage(config.DBConnectionURI, config.DBSchema)
	}

	return &Context{
		Storage:  storage,
		Config:   config,
		Validate: validator.New(),
	}
}

// GetValidationErrorDescription returns a formatted string with first error field, tag and value
func (hc *Context) GetValidationErrorDescription(err error) string {
	return hc.GetValidationErrorDescriptionWithFieldAndValue(err, nil, nil)
}

// GetValidationErrorDescriptionWithField returns a formatted string with given field and first error tag and value
func (hc *Context) GetValidationErrorDescriptionWithField(err error, fieldName string) string {
	return hc.GetValidationErrorDescriptionWithFieldAndValue(err, &fieldName, nil)
}

// GetValidationErrorDescriptionWithValue returns a formatted string with given value and first error tag and field
func (hc *Context) GetValidationErrorDescriptionWithValue(err error, value string) string {
	return hc.GetValidationErrorDescriptionWithFieldAndValue(err, nil, &value)
}

// GetValidationErrorDescriptionWithFieldAndValue returns a formatted string with given field and value and first error tag
func (hc *Context) GetValidationErrorDescriptionWithFieldAndValue(err error, givenFieldName, givenValue *string) string {
	var value interface{}
	fieldName := ""

	validationErrors := err.(validator.ValidationErrors)
	firstError := validationErrors[0]

	if givenFieldName == nil {
		fieldName = firstError.Field()
	} else {
		fieldName = *givenFieldName
	}

	if givenValue == nil {
		value = firstError.Value()
	} else {
		value = *givenValue
	}

	return fmt.Sprintf(errorcode.InvalidDescription, fieldName, firstError.Tag(), value)
}

// LogAndSendError logs and sends the error
func (hc *Context) LogAndSendError(c *gin.Context, causeError error, code, description string, status int) {
	if causeError != nil {
		log.WithError(causeError).Error(code)
	}

	err := &model.Error{
		Code:        code,
		Description: description,
	}

	c.AbortWithStatusJSON(status, err)
}

// UnmarshalPayload unmarshal request payload from context into object
func (hc *Context) UnmarshalPayload(c *gin.Context, i interface{}) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(body), i)
	if err != nil {
		return err
	}
	return nil
}
