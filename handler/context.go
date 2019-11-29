package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/Sehsyha/crounch-back/configuration"
	"github.com/Sehsyha/crounch-back/errorcode"
	"github.com/Sehsyha/crounch-back/model"
	"github.com/Sehsyha/crounch-back/storage"
	"github.com/Sehsyha/crounch-back/storage/mock"
	"github.com/Sehsyha/crounch-back/storage/postgres"
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

// GetValidationErrorDescription returns a formatted string with first error field and tag
func (hc *Context) GetValidationErrorDescription(err error) string {
	validationErrors := err.(validator.ValidationErrors)
	firstError := validationErrors[0]
	return fmt.Sprintf(errorcode.InvalidDescription, firstError.Field(), firstError.Tag())
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
