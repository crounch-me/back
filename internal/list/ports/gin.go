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

func (h *GinServer) GetUserLists(c *gin.Context) {
	userUUID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	lists, err := h.listService.GetUserLists(userUUID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	listsResponse := make([]*List, 0)
	for _, list := range lists {
		products := make([]*Product, 0)
		for _, p := range list.Products() {
			product := &Product{
				UUID: p.UUID(),
			}
			products = append(products, product)
		}

		listResponse := &List{
			UUID:         list.UUID(),
			Name:         list.Name(),
			CreationDate: list.CreationDate(),
			Contributors: list.Contributors(),
			Products:     products,
		}

		listsResponse = append(listsResponse, listResponse)
	}

	response := utils.NewDataResponse(listsResponse)
	c.JSON(http.StatusOK, response)
}
