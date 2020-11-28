package ports

import (
	"net/http"

	"github.com/crounch-me/back/internal/common/utils"
	"github.com/crounch-me/back/internal/list/app"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	listService *app.ListService
}

func NewGinServer(listService *app.ListService) *GinServer {
	return &GinServer{listService: listService}
}

func (h *GinServer) CreateList(c *gin.Context) {
	list := &CreateListRequest{}

	err := utils.UnmarshalPayload(c.Request.Body, list)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userUUID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	listUUID, err := h.listService.CreateList(userUUID, list.Name)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header(utils.HeaderContentLocation, "/lists/"+listUUID)
	c.Status(http.StatusNoContent)
}
