package handler

import (
	"net/http"

	"github.com/Sehsyha/crounch-back/errorcode"
	"github.com/Sehsyha/crounch-back/model"
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

func (hc *Context) AddOFFProduct(c *gin.Context) {
	userID, exists := c.Get(ContextUserID)

	if !exists {
		hc.LogAndSendError(c, nil, errorcode.UserDataCode, errorcode.UserDataDescription, http.StatusInternalServerError)
		return
	}

	id := c.Param("id")

	offProduct := &model.OFFProduct{}

	err := hc.UnmarshalPayload(c, offProduct)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.UnmarshalCode, errorcode.UnmarshalDescription, http.StatusBadRequest)
		return
	}

	err = hc.Validate.Struct(offProduct)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.InvalidCode, hc.GetValidationErrorDescription(err), http.StatusBadRequest)
		return
	}

	l, err := hc.Storage.GetList(id)

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

	if l.Owner.ID != userID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	err = hc.Storage.AddOFFProductToList(id, offProduct)

	if err != nil {
		hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, offProduct)
}
