package handler

import (
	"net/http"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/products"
	"github.com/gin-gonic/gin"
)

const (
	ProductNotFoundDescription = "The product was not found"
)

func (hc *Context) CreateProduct(c *gin.Context) {
	product := &products.Product{}

	err := hc.UnmarshalAndValidate(c, product)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	userID, exists := c.Get(ContextUserID)
	if !exists {
		hc.LogAndSendError(c, domain.NewError(domain.UnknownErrorCode))
		return
	}

	product, err = hc.Services.Product.CreateProduct(product.Name, userID.(string))

	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, product)
}
