package handler

import (
	"net/http"

	"github.com/Sehsyha/crounch-back/errorcode"
	"github.com/Sehsyha/crounch-back/model"
	"github.com/gin-gonic/gin"
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
