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
// @Summary Return health of application
// @ID get-health
// @Produce  json
// @Success 200 {object} domain.Health
// @Failure 500 "Internal Server Error"
// @Router /health [get]
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
