package handler

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sehsyha/crounch-back/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var version string

// Health handles to response to health endpoint (version and alive)
func (hc *Context) Health(c *gin.Context) {
	if version == "" {
		versionBytes, err := ioutil.ReadFile("VERSION")
		if err != nil {
			log.Error(err)
		}
		version = string(versionBytes)
		version = strings.TrimSuffix(version, "\n")
	}
	health := &model.Health{}
	health.Alive = true
	health.Version = version
	c.JSON(http.StatusOK, health)
}
