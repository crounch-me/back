package handler

import (
	"net/http"

	list "github.com/crounch-me/back/internal/listing"
	"github.com/crounch-me/back/util"
	"github.com/gin-gonic/gin"
)

const (
	ListNotFoundDescription = "The list was not found"
)

// AddProductToList handles the request to add a product to a list
// @Summary Add the product to the list
// @ID add-product-to-list
// @Tags product-in-list
// @Produce json
// @Param listID path string true "List ID"
// @Param productID path string true "Product ID"
// @Success 200 {object} list.ProductInListLink
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /lists/{listID}/products/{productID} [post]
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
// @Summary Update the product in the list partially
// @ID update-product-in-list
// @Tags product-in-list
// @Accept json
// @Produce json
// @Param productInList body list.UpdateProductInList true "Product in list"
// @Param listID path string true "Product in list"
// @Param productID path string true "Product in list"
// @Success 200 {object} list.ProductInListLink
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /lists/{listID}/products/{productID} [patch]
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

	updateProductInList := &list.UpdateProductInList{}

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
// @Summary Delete the product from the list
// @ID delete-product-from-list
// @Tags product-in-list
// @Produce json
// @Param listID path string true "List ID"
// @Param productID path string true "Product ID"
// @Success 204
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /lists/{listID}/products/{productID} [delete]
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

// DeleteList deletes a list and all its product links
// @Summary Delete the entire list with its products
// @ID delete-list
// @Tags list
// @Produce json
// @Param listID path string true "List ID"
// @Success 204
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /lists/{listID} [delete]
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
		hc.LogAndSendError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ArchiveList archives a list, it will be marked as readonly
// @Summary Setup archivation date for the list
// @ID update-list
// @Tag list
// @Produce json
// @Param listID path string true "List ID"
// @Success 200 {object} builders.GetListResponse
// @Failure 400 {object} errors.Error
// @Failure 404 {object} errors.Error
// @Failure 500 {object} errors.Error
func (hc *Context) ArchiveList(c *gin.Context) {
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

	list, err := hc.Services.List.ArchiveList(listID, userID)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	listResponse := hc.Builders.List.GetList(list)

	c.JSON(http.StatusOK, listResponse)
}
