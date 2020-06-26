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

type ProductSearchRequest struct {
	Name string `json:"name"`
}

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

func (hc *Context) SearchDefaultProducts(c *gin.Context) {
	productSearchRequest := &ProductSearchRequest{}

	err := hc.UnmarshalAndValidate(c, productSearchRequest)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	products, err := hc.Services.Product.SearchDefaults(productSearchRequest.Name)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.JSON(http.StatusOK, products)
}
