package handler

import (
	"net/http"

	"github.com/Sehsyha/crounch-back/errorcode"
	"github.com/Sehsyha/crounch-back/model"
	"github.com/gin-gonic/gin"
)

func (hc *Context) CreateProduct(c *gin.Context) {
	product := &model.Product{}

	err := hc.UnmarshalPayload(c, product)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.UnmarshalCode, errorcode.UnmarshalDescription, http.StatusBadRequest)
		return
	}

	err = hc.Validate.Struct(product)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.InvalidCode, hc.GetValidationErrorDescription(err), http.StatusBadRequest)
		return
	}

	userID, exists := c.Get(ContextUserID)

	if !exists {
		hc.LogAndSendError(c, nil, errorcode.UserDataCode, errorcode.UserDataDescription, http.StatusInternalServerError)
		return
	}

	product.Owner = &model.User{
		ID: userID.(string),
	}

	err = hc.Storage.CreateProduct(product)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, product)
}
