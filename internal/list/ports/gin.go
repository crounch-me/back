package ports

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/common/utils"
	"github.com/crounch-me/back/internal/list/app"
	"github.com/crounch-me/back/util"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

type GinServer struct {
	listService *app.ListService
	validator   *util.Validator
}

func NewGinServer(listService *app.ListService, validator *util.Validator) (*GinServer, error) {
	if listService == nil {
		return nil, errors.New("listService is nil")
	}

	if validator == nil {
		return nil, errors.New("validator is nil")
	}

	return &GinServer{
		listService: listService,
		validator:   validator,
	}, nil
}

func (h *GinServer) CreateList(c *gin.Context) {
	list := &CreateListRequest{}

	err := utils.UnmarshalPayload(c.Request.Body, list)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.UnmarshalErrorCode))
		return
	}

	err = h.validator.Struct(list)
	if err != nil {
		fields := make([]*internal.FieldError, 0)
		for _, e := range err.(validator.ValidationErrors) {
			field := &internal.FieldError{
				Error: e.Tag(),
				Name:  e.Field(),
			}
			fields = append(fields, field)
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.InvalidErrorCode).WithFields(fields))
		return
	}

	userUUID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, internal.NewError(internal.ForbiddenErrorCode))
		return
	}

	listUUID, err := h.listService.CreateList(userUUID, list.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, internal.NewError(internal.UnknownErrorCode))
		return
	}

	c.Header(utils.HeaderContentLocation, "/lists/"+listUUID)
	c.Status(http.StatusCreated)
}

func (h *GinServer) GetUserLists(c *gin.Context) {
	userUUID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	lists, err := h.listService.GetUserLists(userUUID)
	if err != nil {
		fmt.Println(err)
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
