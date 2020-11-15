package handler

import (
	"net/http"

	"github.com/crounch-me/back/domain"
	"github.com/gin-gonic/gin"
)

const (
	ProductNotFoundDescription = "The product was not found"
)

type CreateProductRequest struct {
  Name string `json:"name" validate:"required"`
}

// CreateProduct creates a new product for the user
// @Summary Create a new product
// @ID create-product
// @Tags product
// @Produce json
// @Param product body CreateProductRequest true "Product to create"
// @Success 200 {object} products.Product
// @Failure 500 {object} domain.Error
// @Security ApiKeyAuth
// @Router /products [post]
func (hc *Context) CreateProduct(c *gin.Context) {
	product := &CreateProductRequest{}

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

	createdProduct, err := hc.Services.Product.CreateProduct(product.Name, userID.(string))

	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdProduct)
}

type ProductSearchRequest struct {
	Name string `json:"name" validate:"gt=2,lt=16"`
}

// SearchDefaultProducts search and return default products
// @Summary Search a product by its name in default products, it removes accentuated characters and is case insensitive
// @ID search-default-products
// @Tags product
// @Produce json
// @Param product body ProductSearchRequest true "Product search request"
// @Success 200 {object} []products.Product
// @Failure 500 {object} domain.Error
// @Security ApiKeyAuth
// @Router /products/search [post]
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
