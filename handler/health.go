package handler

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/crounch-me/back/domain"
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
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		version = strings.TrimSuffix(string(versionBytes), "\n")
	}

	health := &domain.Health{
		Alive:   true,
		Version: version,
	}

	c.JSON(http.StatusOK, health)
}
