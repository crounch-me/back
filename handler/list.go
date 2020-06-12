package handler

import (
	"net/http"

	"github.com/crounch-me/back/errorcode"
	"github.com/crounch-me/back/model"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

const (
	ListNotFoundDescription = "The list was not found"
)

func (hc *Context) CreateList(c *gin.Context) {
	list := &model.List{}

	err := hc.UnmarshalPayload(c, list)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.UnmarshalCode, errorcode.UnmarshalDescription, http.StatusBadRequest)
		return
	}

	err = hc.Validate.Struct(list)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.InvalidCode, hc.GetValidationErrorDescription(err), http.StatusBadRequest)
		return
	}

	userID, exists := c.Get(ContextUserID)

	if !exists {
		hc.LogAndSendError(c, nil, errorcode.UserDataCode, errorcode.UserDataDescription, http.StatusInternalServerError)
		return
	}

	list.Owner = &model.User{
		ID: userID.(string),
	}

	err = hc.Storage.CreateList(list)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, list)
}

func (hc *Context) GetOwnerLists(c *gin.Context) {
	userID, exists := c.Get(ContextUserID)

	if !exists {
		hc.LogAndSendError(c, nil, errorcode.UserDataCode, errorcode.UserDataDescription, http.StatusInternalServerError)
		return
	}

	lists, err := hc.Storage.GetOwnerLists(userID.(string))

	if err != nil {
		hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
		return
	}

	log.WithField("lists", lists).Debug("Response: lists")
	c.JSON(http.StatusOK, lists)
}

func (hc *Context) AddProductToList(c *gin.Context) {
	userID, exists := c.Get(ContextUserID)

	if !exists {
		hc.LogAndSendError(c, nil, errorcode.UserDataCode, errorcode.UserDataDescription, http.StatusInternalServerError)
		return
	}

	listID := c.Param("listID")

	err := hc.Validate.Var(listID, "uuid")
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.InvalidCode, hc.GetValidationErrorDescriptionWithField(err, "list ID"), http.StatusBadRequest)
		return
	}

	productID := c.Param("productID")

	err = hc.Validate.Var(productID, "uuid")
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.InvalidCode, hc.GetValidationErrorDescriptionWithField(err, "product ID"), http.StatusBadRequest)
		return
	}

	list, err := hc.Storage.GetList(listID)

	if err != nil {
		if databaseError, ok := err.(*model.DatabaseError); ok {
			switch databaseError.Type {
			case model.ErrNotFound:
				hc.LogAndSendError(c, err, errorcode.NotFoundCode, ListNotFoundDescription, http.StatusNotFound)
				return
			}
		}
		hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
		return
	}

	if list.Owner.ID != userID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	product, err := hc.Storage.GetProduct(productID)

	if err != nil {
		if databaseError, ok := err.(*model.DatabaseError); ok {
			switch databaseError.Type {
			case model.ErrNotFound:
				hc.LogAndSendError(c, err, errorcode.NotFoundCode, ProductNotFoundDescription, http.StatusNotFound)
				return
			}
		}
		hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
		return
	}

	if product.Owner.ID != userID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	productInList, err := hc.Storage.GetProductInList(productID, listID)

	if err != nil {
		if _, ok := err.(*model.DatabaseError); ok {
			// No return when not found, it's normal because the association must not exist
		} else {
			hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
			return
		}
	}

	if productInList != nil {
		hc.LogAndSendError(c, err, errorcode.DuplicateCode, errorcode.DuplicateDescription, http.StatusConflict)
		return
	}

	err = hc.Storage.AddProductToList(productID, listID)

	if err != nil {
		hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
		return
	}

	productInList = &model.ProductInList{
		ProductID: productID,
		ListID:    listID,
	}

	c.JSON(http.StatusCreated, productInList)
}
