package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/crounch-me/back/builders"
	"github.com/crounch-me/back/configuration"
	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/list"
	"github.com/crounch-me/back/internal/products"
	"github.com/crounch-me/back/storage"
	"github.com/crounch-me/back/storage/postgres"
	"github.com/crounch-me/back/util"
)

const (
	// ContextUserID is the key to retrieve user id from gin context
	ContextUserID = "ContextUserID"
)

type Services struct {
	List    *list.ListService
	Product *products.ProductService
}

type Builders struct {
	List *builders.ListBuilder
}

// Context holds everything to respond to requests
type Context struct {
	Generation internal.Generation
	Storage    storage.Storage
	Config     *configuration.Config
	Validator  *util.Validator
	Services   *Services
	Builders   *Builders
}

// NewContext creates and initialize everything for the requests
func NewContext(config *configuration.Config) *Context {
	var storage storage.Storage

	postgres.InitDB(config.DBConnectionURI)
	storage = postgres.NewStorage(config.DBConnectionURI, config.DBSchema)
	generation := &util.GenerationImpl{}

	validator := util.NewValidator()

	return &Context{
		Storage:   storage,
		Config:    config,
		Validator: validator,
		Services:  NewServices(storage, generation),
		Builders:  NewBuilders(),
	}
}

// LogAndSendError logs and sends the error
func (hc *Context) LogAndSendError(c *gin.Context, err error) {
	var status int
	if internalErr, ok := err.(*internal.Error); ok {
		status = hc.ErrorCodeToHTTPStatus(internalErr.Code)
		logBuilder := logrus.WithTime(time.Now())

		if internalErr.Call != nil {
			logBuilder = logBuilder.
				WithField("method", internalErr.Call.MethodName).
				WithField("package", internalErr.Call.PackageName)
		}

		if internalErr.Cause != nil {
			logBuilder = logBuilder.WithError(internalErr.Cause)
		}

		if status < http.StatusBadRequest || status >= http.StatusInternalServerError {
			logBuilder.Error(internalErr.Code)
		} else {
			if internalErr.Fields != nil {
				for _, field := range internalErr.Fields {
					logBuilder = logBuilder.WithField(field.Name, field.Error)
				}
			}

			logBuilder.Debug(internalErr.Code)
		}
	} else {
		status = http.StatusInternalServerError
		logrus.Error(err)
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

func (hc *Context) GetUserIDFromContext(c *gin.Context) (string, *internal.Error) {
	userID, exists := c.Get(ContextUserID)
	if !exists {
		return "", internal.NewError(internal.UnknownErrorCode).WithCall("handler", "GetUserIDFromContext")
	}
	return userID.(string), nil
}

func (hc *Context) UnmarshalAndValidate(c *gin.Context, i interface{}) *internal.Error {
	err := hc.UnmarshalPayload(c, i)

	if err != nil {
		return internal.NewError(internal.UnmarshalErrorCode).WithCause(err)
	}

	err = hc.Validator.Struct(i)
	if err != nil {
		fields := make([]*internal.FieldError, 0)
		for _, e := range err.(validator.ValidationErrors) {
			field := &internal.FieldError{
				Error: e.Tag(),
				Name:  e.Field(),
			}
			fields = append(fields, field)
		}
		return internal.NewError(internal.InvalidErrorCode).WithFields(fields)
	}

	return nil
}

// NewBuilders create an object which contains all necessary builders
func NewBuilders() *Builders {
	return &Builders{
		List: &builders.ListBuilder{},
	}
}

// NewServices create an object which contains all necessary services
func NewServices(storage storage.Storage, generation internal.Generation) *Services {
	return &Services{
		List: &list.ListService{
			ListStorage:        storage,
			ProductStorage:     storage,
			ContributorStorage: storage,
			Generation:         generation,
		},
		Product: &products.ProductService{
			ProductStorage: storage,
			Generation:     generation,
		},
	}
}
