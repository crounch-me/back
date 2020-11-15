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
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/authorization.go"
	"github.com/crounch-me/back/domain/lists"
	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
	"github.com/crounch-me/back/storage"
	"github.com/crounch-me/back/storage/postgres"
	"github.com/crounch-me/back/util"
)

const (
	// ContextUserID is the key to retrieve user id from gin context
	ContextUserID = "ContextUserID"
)

type Services struct {
	Authorization *authorization.AuthorizationService
	List          *lists.ListService
	Product       *products.ProductService
	User          *users.UserService
}

type Builders struct {
  List *builders.ListBuilder
}

// Context holds everything to respond to requests
type Context struct {
	Generation domain.Generation
	Storage    storage.Storage
	Config     *configuration.Config
	Validator  *util.Validator
  Services   *Services
  Builders *Builders
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
    Builders: NewBuilders(),
	}
}

// LogAndSendError logs and sends the error
func (hc *Context) LogAndSendError(c *gin.Context, err error) {
	var status int
	if domainErr, ok := err.(*domain.Error); ok {
		status = hc.ErrorCodeToHTTPStatus(domainErr.Code)
		logBuilder := logrus.WithTime(time.Now())

		if domainErr.Call != nil {
			logBuilder = logBuilder.
				WithField("method", domainErr.Call.MethodName).
				WithField("package", domainErr.Call.PackageName)
		}

		if domainErr.Cause != nil {
			logBuilder = logBuilder.WithError(domainErr.Cause)
		}

		if status < http.StatusBadRequest || status >= http.StatusInternalServerError {
			logBuilder.Error(domainErr.Code)
		} else {
			if domainErr.Fields != nil {
				for _, field := range domainErr.Fields {
					logBuilder = logBuilder.WithField(field.Name, field.Error)
				}
			}

			logBuilder.Debug(domainErr.Code)
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

func (hc *Context) GetUserIDFromContext(c *gin.Context) (string, *domain.Error) {
	userID, exists := c.Get(ContextUserID)
	if !exists {
		return "", domain.NewError(domain.UnknownErrorCode).WithCall("handler", "GetUserIDFromContext")
	}
	return userID.(string), nil
}

func (hc *Context) UnmarshalAndValidate(c *gin.Context, i interface{}) *domain.Error {
	err := hc.UnmarshalPayload(c, i)

	if err != nil {
		return domain.NewError(domain.UnmarshalErrorCode).WithCause(err)
	}

	err = hc.Validator.Struct(i)
	if err != nil {
		fields := make([]*domain.FieldError, 0)
		for _, e := range err.(validator.ValidationErrors) {
			field := &domain.FieldError{
				Error: e.Tag(),
				Name:  e.Field(),
			}
			fields = append(fields, field)
		}
		return domain.NewError(domain.InvalidErrorCode).WithFields(fields)
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
func NewServices(storage storage.Storage, generation domain.Generation) *Services {
	return &Services{
		Authorization: &authorization.AuthorizationService{
			AuthorizationStorage: storage,
			UserStorage:          storage,
			Generation:           generation,
		},
		List: &lists.ListService{
			ListStorage:    storage,
      ProductStorage: storage,
      ContributorStorage: storage,
			Generation:     generation,
		},
		Product: &products.ProductService{
			ProductStorage: storage,
			Generation:     generation,
		},
		User: &users.UserService{
			UserStorage: storage,
			Generation:  generation,
		},
	}
}
