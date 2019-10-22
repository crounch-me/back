package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/Sehsyha/crounch-back/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var version string

func (hc *Context) Health(c *gin.Context) {
	if version == "" {
		versionBytes, err := ioutil.ReadFile("VERSION")
		if err != nil {
			log.Error(err)
		}
		version = string(versionBytes)
	}
	health := &model.Health{}
	health.Alive = true
	health.Version = version
	c.JSON(http.StatusOK, health)
}
