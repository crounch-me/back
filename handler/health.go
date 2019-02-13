package handler

import (
	"net/http"

	"github.com/Sehsyha/crounch-back/model"
	"github.com/gin-gonic/gin"
)

func (hc *Context) Health(c *gin.Context) {
	health := &model.Health{}
	health.Alive = true
	c.JSON(http.StatusOK, health)
}
