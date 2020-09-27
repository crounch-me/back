package handler

import (
	"net/http"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/lists"
	"github.com/crounch-me/back/util"
	"github.com/gin-gonic/gin"
)

const (
	ListNotFoundDescription = "The list was not found"
)

func (hc *Context) CreateList(c *gin.Context) {
	list := &lists.List{}

	err := hc.UnmarshalAndValidate(c, list)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	userID, exists := c.Get(ContextUserID)
	if !exists {
		hc.LogAndSendError(c, domain.NewError(domain.UnknownErrorCode))
		return
	}

	list, err = hc.Services.List.CreateList(list.Name, userID.(string))
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, list)
}

// GetOwnerLists handles the request to get the owner's lists
func (hc *Context) GetOwnerLists(c *gin.Context) {
	userID, exists := c.Get(ContextUserID)
	if !exists {
		hc.LogAndSendError(c, domain.NewError(domain.UnknownErrorCode))
		return
	}

	lists, err := hc.Services.List.GetOwnersLists(userID.(string))
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, lists)
}

// AddProductToList handles the request to add a product to a list
func (hc *Context) AddProductToList(c *gin.Context) {
	userID, err := hc.GetUserIDFromContext(c)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	listID := c.Param("listID")
	err = hc.Validator.Var("listID", listID, "uuid")
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	productID := c.Param("productID")
	err = hc.Validator.Var("productID", productID, "uuid")
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	productInList, err := hc.Services.List.AddProductToList(productID, listID, userID)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, productInList)
}

// UpdateProductInList updates the product in list partially
func (hc *Context) UpdateProductInList(c *gin.Context) {
	logger := util.GetLogger()
	userID, err := hc.GetUserIDFromContext(c)

	if err != nil {
		logger.WithError(err).Debug("error during user id retrieving")
		hc.LogAndSendError(c, err)
		return
	}

	listID := c.Param("listID")
	err = hc.Validator.Var("listID", listID, "uuid")
	if err != nil {
		logger.WithError(err).Debug("error during list id validation")
		hc.LogAndSendError(c, err)
		return
	}

	productID := c.Param("productID")
	err = hc.Validator.Var("productID", productID, "uuid")
	if err != nil {
		logger.WithError(err).Debug("error during product id validation")
		hc.LogAndSendError(c, err)
		return
	}

	updateProductInList := &lists.UpdateProductInList{}

	err = hc.UnmarshalAndValidate(c, updateProductInList)
	if err != nil {
		logger.WithError(err).Debug("error during product in list body validation")
		hc.LogAndSendError(c, err)
		return
	}

	productInListLink, err := hc.Services.List.UpdateProductInList(updateProductInList, productID, listID, userID)
	if err != nil {
		logger.WithError(err).Debug("error during product in list link update")
		hc.LogAndSendError(c, err)
		return
	}

	logger.WithField("productInListLink", productInListLink).
		Debug("OK - Request succeeded")

	c.JSON(http.StatusOK, productInListLink)
}

// DeleteProductFromList removes the product from the list
func (hc *Context) DeleteProductFromList(c *gin.Context) {
	userID, err := hc.GetUserIDFromContext(c)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	listID := c.Param("listID")
	err = hc.Validator.Var("listID", listID, "uuid")
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	productID := c.Param("productID")
	err = hc.Validator.Var("productID", productID, "uuid")
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	err = hc.Services.List.DeleteProductFromList(productID, listID, userID)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (hc *Context) DeleteList(c *gin.Context) {
	logger := util.GetLogger()
	userID, err := hc.GetUserIDFromContext(c)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	listID := c.Param("listID")
	err = hc.Validator.Var("listID", listID, "uuid")
	if err != nil {
		logger.WithError(err).Debug("error during list id validation")
		hc.LogAndSendError(c, err)
		return
	}

	err = hc.Services.List.DeleteList(listID, userID)
	if err != nil {
		logger.WithError(err).Debug("error during list deletion")
		hc.LogAndSendError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (hc *Context) GetList(c *gin.Context) {
	userID, err := hc.GetUserIDFromContext(c)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	listID := c.Param("listID")
	err = hc.Validator.Var("listID", listID, "uuid")
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	list, err := hc.Services.List.GetList(listID, userID)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, list)
}
