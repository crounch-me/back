package handler

import (
	"net/http"

	"github.com/crounch-me/back/domain"
	"github.com/gin-gonic/gin"
)

// Health handles to response to health endpoint and tells if service is alive or not
// @Summary Return health of application
// @ID get-health
// @Produce  json
// @Success 200 {object} domain.Health
// @Failure 500 "Internal Server Error"
// @Router /health [get]
func (hc *Context) Health(c *gin.Context) {
	health := &domain.Health{
		Alive:   true,
	}

	c.JSON(http.StatusOK, health)
}
