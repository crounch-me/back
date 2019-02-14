package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Sehsyha/crounch-back/configuration"
	"github.com/Sehsyha/crounch-back/errorcode"
	"github.com/Sehsyha/crounch-back/model"
	"github.com/Sehsyha/crounch-back/storage"
	"github.com/Sehsyha/crounch-back/storage/mock"
	"github.com/Sehsyha/crounch-back/storage/neo"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

type Context struct {
	Storage  storage.Storage
	Config   *configuration.Config
	Validate *validator.Validate
}

func NewContext(config *configuration.Config) *Context {
	var storage storage.Storage

	if config.Mock {
		storage = mock.NewStorageMock()
	} else {
		storage = neo.NewNeoStorage()
	}

	return &Context{
		Storage:  storage,
		Config:   config,
		Validate: validator.New(),
	}
}

func (hc *Context) GetValidationErrorDescription(err error) string {
	validationErrors := err.(validator.ValidationErrors)
	firstError := validationErrors[0]
	return fmt.Sprintf(errorcode.InvalidDescription, firstError.Field(), firstError.Tag())
}

func (hc *Context) LogAndSendError(c *gin.Context, causeError error /*, title*/, code, description string, status int) {
	if causeError != nil {
		log.Error(causeError)
	}

	err := &model.Error{
		Code:        code,
		Description: description,
	}

	c.AbortWithStatusJSON(status, err)
}

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
