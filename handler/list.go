package handler

import (
	"net/http"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/lists"
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

	lists, err := hc.Services.List.GetOwnerLists(userID.(string))
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, lists)
}

// AddProductToList handles the request to add a product to a list
func (hc *Context) AddProductToList(c *gin.Context) {
	userID, domainErr := hc.GetUserIDFromContext(c)
	if domainErr != nil {
		hc.LogAndSendError(c, domainErr)
		return
	}

	listID := c.Param("listID")

	err := hc.Validate.Var(listID, "uuid")
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	productID := c.Param("productID")

	err = hc.Validate.Var(productID, "uuid")
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	productInList, domainErr := hc.Services.List.AddProductToList(productID, listID, userID)
	if domainErr != nil {
		hc.LogAndSendError(c, domainErr)
		return
	}

	c.JSON(http.StatusCreated, productInList)
}
