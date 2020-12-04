package ports

import (
	"errors"
	"fmt"
	"net/http"

	accountApp "github.com/crounch-me/back/internal/account/app"
	commonErrors "github.com/crounch-me/back/internal/common/errors"
	"github.com/crounch-me/back/internal/common/server"
	listApp "github.com/crounch-me/back/internal/listing/app"
	"github.com/crounch-me/back/util"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

const (
	listUUIDPathParam = "listID"
	listPath          = "/listing/lists"
	listWithIDPath    = "/listing/lists/:listID"
	archiveListPath   = "/listing/lists/:listID/archive"
)

type GinServer struct {
	accountService *accountApp.AccountService
	listService    *listApp.ListService
	validator      *util.Validator
}

func NewGinServer(listService *listApp.ListService, accountService *accountApp.AccountService, validator *util.Validator) (*GinServer, error) {
	if listService == nil {
		return nil, errors.New("listService is nil")
	}

	if accountService == nil {
		return nil, errors.New("accountService is nil")
	}

	if validator == nil {
		return nil, errors.New("validator is nil")
	}

	return &GinServer{
		listService: listService,
		validator:   validator,
	}, nil
}

func (s *GinServer) ConfigureRoutes(r *gin.Engine) {
	r.POST(listPath, server.CheckUserAuthorization(s.accountService), s.CreateList)
	r.GET(listPath, server.CheckUserAuthorization(s.accountService), s.GetContributorsLists)
	r.OPTIONS(listPath, server.OptionsHandler([]string{http.MethodGet, http.MethodPost}))
}

// CreateList creates a new list
// @Summary Create a list
// @ID create-list
// @Tags listing
// @Accept json
// @Produce  json
// @Param list body CreateListRequest true "List to create"
// @Success 201
// @Failure 400 {object} errors.Error
// @Failure 403 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /listing/lists [post]
func (s *GinServer) CreateList(c *gin.Context) {
	list := &CreateListRequest{}

	err := server.UnmarshalPayload(c.Request.Body, list)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, commonErrors.NewError(commonErrors.UnmarshalErrorCode))
		return
	}

	err = s.validator.Struct(list)
	if err != nil {
		fields := make([]*commonErrors.FieldError, 0)
		for _, e := range err.(validator.ValidationErrors) {
			field := &commonErrors.FieldError{
				Error: e.Tag(),
				Name:  e.Field(),
			}
			fields = append(fields, field)
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, commonErrors.NewError(commonErrors.InvalidErrorCode).WithFields(fields))
		return
	}

	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, commonErrors.NewError(commonErrors.ForbiddenErrorCode))
		return
	}

	listUUID, err := s.listService.CreateList(userUUID, list.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonErrors.NewError(commonErrors.UnknownErrorCode))
		return
	}

	c.Header(server.HeaderContentLocation, "/lists/"+listUUID)
	c.Status(http.StatusCreated)
}

// GetContributorsLists get the authenticated contributor accessible lists
// @Summary Get the authenticated contributor accessible lists
// @ID get-contributors-lists
// @Tags listing
// @Produce  json
// @Success 200 {object} []List
// @Failure 403 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /listing/lists [get]
func (s *GinServer) GetContributorsLists(c *gin.Context) {
	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	lists, err := s.listService.GetContributorLists(userUUID)
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

		contributors := make([]*Contributor, 0)
		for _, c := range list.Contributors() {
			contributor := &Contributor{
				UUID: c.UUID(),
			}
			contributors = append(contributors, contributor)
		}

		listResponse := &List{
			UUID:         list.UUID(),
			Name:         list.Name(),
			CreationDate: list.CreationDate(),
			Contributors: contributors,
			Products:     products,
		}

		listsResponse = append(listsResponse, listResponse)
	}

	server.JSON(c, listsResponse)
}

// GetList reads a specific list with its product and contributor ids
// @Summary Reads a list with its product and contributor ids
// @ID get-list
// @Tags list
// @Produce json
// @Param listID path string true "List ID"
// @Success 200 {object} builders.GetListResponse
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /lists/{listID} [get]
func (s *GinServer) GetList(c *gin.Context) {
	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	listUUID := c.Param(listUUIDPathParam)
	list, err := s.listService.ReadList(userUUID, listUUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonErrors.NewError(commonErrors.UnknownErrorCode))
		return
	}

	server.JSON(c, list)
}
