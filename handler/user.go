package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Sehsyha/crounch-back/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (hc *Context) Signup(c *gin.Context) {
	u := &model.User{}

	if UnmarshalPayload(c, u) {
		return
	}

	err := hc.Storage.CreateUser(u)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(500)
		return
	}

	c.JSON(http.StatusCreated, u)
}

func (hc *Context) Login(c *gin.Context) {
	u := &model.User{}

	if UnmarshalPayload(c, u) {
		return
	}

	authorization, err := hc.Storage.CreateAuthorization(u)

	if err != nil {
		log.Error(err)
		c.AbortWithStatus(500)
		return
	}

	c.JSON(http.StatusCreated, authorization)
}

func UnmarshalPayload(c *gin.Context, i interface{}) bool {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "unmarshall error",
		})
		return true
	}

	err = json.Unmarshal([]byte(body), i)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "unmarshall error",
		})
		return true
	}
	return false
}
