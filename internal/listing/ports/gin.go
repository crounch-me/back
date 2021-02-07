package ports

import (
	"errors"
	"net/http"

	accountApp "github.com/crounch-me/back/internal/account/app"
	commonErrors "github.com/crounch-me/back/internal/common/errors"
	"github.com/crounch-me/back/internal/common/server"
	"github.com/crounch-me/back/internal/common/utils"
	listApp "github.com/crounch-me/back/internal/listing/app"
	"github.com/crounch-me/back/internal/listing/domain/lists"
	"github.com/crounch-me/back/internal/products/domain/products"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

const (
	listUUIDPathParam = "listID"

	listPath                       = "/listing/lists"
	listWithIDPath                 = "/listing/lists/:listID"
	archiveListPath                = "/listing/lists/:listID/archive"
	listWithIDAndProductWithIDPath = "/listing/lists/:listID/products/:productID"
	buyProductPath                 = "/listing/lists/:listID/products/:productID/buy"

	productAlreadyInListErrorCode = "product-already-in-list-error"
	listNotFoundErrorCode         = "list-not-found-error"
	productNotFoundErrorCode      = "product-not-found-error"
)

type GinServer struct {
	accountService *accountApp.AccountService
	listService    *listApp.ListService
	validator      *utils.Validator
}

func NewGinServer(
	listService *listApp.ListService,
	accountService *accountApp.AccountService,
	validator *utils.Validator,
) (*GinServer, error) {
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
		listService:    listService,
		accountService: accountService,
		validator:      validator,
	}, nil
}

func (s *GinServer) ConfigureRoutes(r *gin.Engine) {
	r.POST(listPath, server.CheckUserAuthorization(s.accountService), s.CreateList)
	r.GET(listPath, server.CheckUserAuthorization(s.accountService), s.GetContributorsLists)
	r.OPTIONS(listPath, server.OptionsHandler([]string{http.MethodGet, http.MethodPost}))

	r.GET(listWithIDPath, server.CheckUserAuthorization(s.accountService), s.GetList)
	r.DELETE(listWithIDPath, server.CheckUserAuthorization(s.accountService), s.DeleteList)
	r.OPTIONS(listWithIDPath, server.OptionsHandler([]string{http.MethodGet}))

	r.POST(archiveListPath, server.CheckUserAuthorization(s.accountService), s.ArchiveList)
	r.OPTIONS(archiveListPath, server.OptionsHandler([]string{http.MethodPost}))

	r.POST(listWithIDAndProductWithIDPath, server.CheckUserAuthorization(s.accountService), s.AddProductToList)
	r.DELETE(listWithIDAndProductWithIDPath, server.CheckUserAuthorization(s.accountService), s.DeleteProductInList)
	r.OPTIONS(listWithIDAndProductWithIDPath, server.OptionsHandler([]string{http.MethodPost}))

	r.PATCH(buyProductPath, server.CheckUserAuthorization(s.accountService), s.BuyProduct)
	r.OPTIONS(buyProductPath, server.OptionsHandler([]string{http.MethodPatch}))
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

	createdList, err := s.listService.CreateList(userUUID, list.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonErrors.NewError(commonErrors.UnknownErrorCode))
		return
	}

	listResponse := NewResponseBuilder().
		FromDomain(createdList).
		Build()

	server.JSON(c, listResponse)
}

// GetContributorsLists get the authenticated contributor accessible lists
// @Summary Get the authenticated contributor accessible lists
// @ID get-contributors-lists
// @Tags listing
// @Produce  json
// @Success 200 {object} []ListResponse
// @Failure 403 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /listing/lists [get]
func (s *GinServer) GetContributorsLists(c *gin.Context) {
	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		logrus.Debug(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	lists, err := s.listService.GetContributorLists(userUUID)
	if err != nil {
		logrus.Debug(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	listsResponse := make([]*ListResponse, 0)
	responseBuilder := NewResponseBuilder()
	for _, list := range lists {
		listResponse := responseBuilder.FromDomain(list).Build()
		responseBuilder.Reset()

		listsResponse = append(listsResponse, listResponse)
	}

	server.JSON(c, listsResponse)
}

// GetList reads a specific list with its product and contributor ids
// @Summary Reads a list with its product and contributor ids
// @ID get-list
// @Tags listing
// @Produce json
// @Param listID path string true "List UUID"
// @Success 200 {object} ListResponse
// @Failure 403 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /listing/lists/{listID} [get]
func (s *GinServer) GetList(c *gin.Context) {
	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	listID := c.Param(listUUIDPathParam)
	commonErr := s.validator.Var(listUUIDPathParam, listID, "uuid")
	if commonErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, commonErrors.NewError(commonErrors.InvalidErrorCode))
		return
	}

	list, err := s.listService.ReadList(userUUID, listID)
	if err != nil {
		logrus.Debug(err)
		if errors.Is(err, lists.ErrUserNotContributor) {
			c.AbortWithStatusJSON(http.StatusForbidden, commonErrors.NewError(commonErrors.ForbiddenErrorCode))
			return
		}

		if errors.Is(err, lists.ErrListNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, commonErrors.NewError(listNotFoundErrorCode))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, commonErrors.NewError(commonErrors.UnknownErrorCode))
		return
	}

	listResponse := NewResponseBuilder().
		FromDomain(list).
		Build()

	server.JSON(c, listResponse)
}

// ArchiveList and mark it as readonly
// @Summary Archives a list and mark it as readonly
// @ID archive-list
// @Tag listing
// @Produce json
// @Param listID path string true "List UUID"
// @Success 204
// @Failure 403 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Router /listing/lists/:listID/archive [archive]
func (s *GinServer) ArchiveList(c *gin.Context) {
	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	listID := c.Param(listUUIDPathParam)
	commonErr := s.validator.Var(listUUIDPathParam, listID, "uuid")
	if commonErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, commonErrors.NewError(commonErrors.InvalidErrorCode))
		return
	}

	err = s.listService.ArchiveList(userUUID, listID)
	if err != nil {
		if errors.Is(err, lists.ErrListNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, commonErrors.NewError(listNotFoundErrorCode))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonErrors.NewError(commonErrors.UnknownErrorCode))
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteList and all its product links and contributors
// @Summary Delete the entire list with its products links and contributors
// @ID delete-list
// @Tags listing
// @Produce json
// @Param listID path string true "List UUID"
// @Success 204
// @Failure 403 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /listing/lists/{listID} [delete]
func (s *GinServer) DeleteList(c *gin.Context) {
	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, commonErrors.NewError(commonErrors.ForbiddenErrorCode))
		return
	}

	listID := c.Param(listUUIDPathParam)
	commonErr := s.validator.Var(listUUIDPathParam, listID, "uuid")
	if commonErr != nil {
		logrus.WithError(err).WithField("test", err).Debug("error while deleting list")
		c.AbortWithStatusJSON(http.StatusBadRequest, commonErrors.NewError(commonErrors.InvalidErrorCode))
		return
	}

	err = s.listService.DeleteList(userUUID, listID)
	if err != nil {
		logrus.Debug(err)
		if errors.Is(err, lists.ErrListNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, commonErrors.NewError(listNotFoundErrorCode))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonErrors.NewError(commonErrors.UnknownErrorCode))
		return
	}

	c.Status(http.StatusNoContent)
}

// AddProductToList adds a product to a list
// @Summary Add the product to the list
// @ID add-product-to-list
// @Tags listing
// @Produce json
// @Param listID path string true "List ID"
// @Param productID path string true "Product ID"
// @Success 200 {object} AddProductToListRequest
// @Failure 400 {object} errors.Error
// @Failure 403 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /listing/lists/products [post]
func (s *GinServer) AddProductToList(c *gin.Context) {
	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, commonErrors.NewError(commonErrors.ForbiddenErrorCode))
		return
	}

	request := &AddProductToListRequest{}
	err = server.UnmarshalPayload(c.Request.Body, request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, commonErrors.NewError(commonErrors.UnmarshalErrorCode))
		return
	}

	err = s.validator.Struct(request)
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

	err = s.listService.AddProductToList(userUUID, request.ProductUUID, request.ListUUID)
	if err != nil {
		if errors.Is(err, lists.ErrUserNotContributor) {
			c.AbortWithStatusJSON(http.StatusForbidden, commonErrors.NewError(commonErrors.ForbiddenErrorCode))
			return
		}

		if errors.Is(err, lists.ErrProductAlreadyInList) {
			c.AbortWithStatusJSON(http.StatusConflict, commonErrors.NewError(productAlreadyInListErrorCode))
			return
		}

		if errors.Is(err, lists.ErrListNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, commonErrors.NewError(listNotFoundErrorCode))
			return
		}

		if errors.Is(err, products.ErrProductNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, commonErrors.NewError(productNotFoundErrorCode))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, commonErrors.NewError(commonErrors.UnknownErrorCode))
		return
	}

	c.Status(http.StatusCreated)
}

// BuyProduct buys the product in the list
// @Summary Buys the product in the list
// @ID buy-product-in-list
// @Tags listing
// @Accept json
// @Produce json
// @Param productInListRequest body BuyProductInListRequest true "Product in list request"
// @Success 204
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /listing/lists/products [patch]
func (s *GinServer) BuyProduct(c *gin.Context) {
	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, commonErrors.NewError(commonErrors.ForbiddenErrorCode))
		return
	}

	request := &BuyProductInListRequest{}
	err = server.UnmarshalPayload(c.Request.Body, request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, commonErrors.NewError(commonErrors.UnmarshalErrorCode))
		return
	}

	err = s.validator.Struct(request)
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

	err = s.listService.BuyProductInList(userUUID, request.ProductUUID, request.ListUUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonErrors.NewError(commonErrors.UnknownErrorCode))
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteProductInList removes the product int the list
// @Summary Delete the product in the list
// @ID delete-product-from-list
// @Tags listing
// @Param listID path string true "List ID"
// @Param productID path string true "Product ID"
// @Success 204
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /listing/lists/{listID}/products/{productID} [delete]
func (s *GinServer) DeleteProductInList(c *gin.Context) {
	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, commonErrors.NewError(commonErrors.ForbiddenErrorCode))
		return
	}

	listID := c.Param(listUUIDPathParam)
	err = s.validator.Var(listUUIDPathParam, listID, "uuid")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, commonErrors.NewError(commonErrors.InvalidErrorCode))
		return
	}
	productUUID := c.Param("productID")

	err = s.listService.DeleteProductInList(userUUID, productUUID, listID)
	if err != nil {
		logrus.Debug(err)
		if errors.Is(err, lists.ErrListNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, commonErrors.NewError(listNotFoundErrorCode))
			return
		}

		if errors.Is(err, lists.ErrUserNotContributor) {
			c.AbortWithStatusJSON(http.StatusForbidden, commonErrors.NewError(commonErrors.ForbiddenErrorCode))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, commonErrors.NewError(commonErrors.UnknownErrorCode))
		return
	}

	c.Status(http.StatusNoContent)
}
