package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Sehsyha/crounch-back/configuration"
	"github.com/Sehsyha/crounch-back/errorcode"
	"github.com/Sehsyha/crounch-back/model"
	"github.com/Sehsyha/crounch-back/storage"
	"github.com/Sehsyha/crounch-back/storage/mock"
	"github.com/Sehsyha/crounch-back/storage/neo"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Context struct {
	Storage storage.Storage
	Config  *configuration.Config
}

func NewContext(config *configuration.Config) *Context {
	var storage storage.Storage

	if config.Mock {
		storage = mock.NewStorageMock()
	} else {
		storage = neo.NewNeoStorage()
	}

	return &Context{
		Storage: storage,
		Config:  config,
	}
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

func (hc *Context) UnmarshalPayload(c *gin.Context, i interface{}) bool {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.UnmarshalCode, errorcode.UnmarshalDescription, http.StatusBadRequest)
		return true
	}

	err = json.Unmarshal([]byte(body), i)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.UnmarshalCode, errorcode.UnmarshalDescription, http.StatusBadRequest)
		return true
	}
	return false
}
