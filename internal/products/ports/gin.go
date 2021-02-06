package ports

import (
	"errors"
	"net/http"

	accountApp "github.com/crounch-me/back/internal/account/app"
	commonErrors "github.com/crounch-me/back/internal/common/errors"
	"github.com/crounch-me/back/internal/common/server"
	"github.com/crounch-me/back/internal/common/utils"
	"github.com/crounch-me/back/internal/products/app"
	_ "github.com/crounch-me/back/internal/products/domain/products"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

const (
	productPath       = "/products"
	productSearchPath = "/products/search"
)

type GinServer struct {
	accountService *accountApp.AccountService
	productService *app.ProductService
	validator      *utils.Validator
}

func NewGinServer(
	productService *app.ProductService,
	accountService *accountApp.AccountService,
	validator *utils.Validator,
) (*GinServer, error) {
	if productService == nil {
		return nil, errors.New("productService is nil")
	}

	if accountService == nil {
		return nil, errors.New("accountService is nil")
	}

	if validator == nil {
		return nil, errors.New("validator is nil")
	}

	return &GinServer{
		accountService: accountService,
		productService: productService,
		validator:      validator,
	}, nil
}

func (s *GinServer) ConfigureRoutes(r *gin.Engine) {
	r.POST(productPath, server.CheckUserAuthorization(s.accountService), s.CreateProduct)
	r.OPTIONS(productPath, server.OptionsHandler([]string{http.MethodPost}))

	r.POST(productSearchPath, server.CheckUserAuthorization(s.accountService), s.SearchDefaultProducts)
	r.OPTIONS(productSearchPath, server.OptionsHandler([]string{http.MethodPost}))
}

// CreateProduct creates a new product, searchable by its creator
// @Summary Create a new product, searchable by its creator
// @ID create-product
// @Tags products
// @Param product body CreateProductRequest true "Product to create"
// @Success 204
// @Failure 400 {object} errors.Error
// @Failure 403 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /products [post]
func (s *GinServer) CreateProduct(c *gin.Context) {
	product := &CreateProductRequest{}

	err := server.UnmarshalPayload(c.Request.Body, product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, commonErrors.NewError(commonErrors.UnmarshalErrorCode))
		return
	}

	err = s.validator.Struct(product)
	if err != nil {
		fields := make([]*commonErrors.FieldError, 0)
		for _, e := range err.(validator.ValidationErrors) {
			field := &commonErrors.FieldError{
				Error: e.Tag(),
				Name:  e.Field(),
			}
			fields = append(fields, field)
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, commonErrors.NewError(commonErrors.InvalidErrorCode).WithFields(fields))
		return
	}

	productUUID, err := s.productService.CreateProduct(product.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonErrors.NewError(commonErrors.UnknownErrorCode))
		return
	}

	c.Header(server.HeaderContentLocation, "/products/"+productUUID)
	c.Status(http.StatusCreated)
}

type ProductSearchRequest struct {
	Name string `json:"name" validate:"gt=2,lt=16"`
}

// SearchDefaultProducts search and return default products
// @Summary Search a product by its name in default products, it removes accentuated characters and is case insensitive
// @ID search-default-products
// @Tags products
// @Produce json
// @Param product body ProductSearchRequest true "Product search request"
// @Success 200 {object} []products.Product
// @Failure 500 {object} errors.Error
// @Security ApiKeyAuth
// @Router /products/search [post]
func (s *GinServer) SearchDefaultProducts(c *gin.Context) {
	productSearchRequest := &ProductSearchRequest{}

	err := server.UnmarshalPayload(c.Request.Body, productSearchRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, commonErrors.NewError(commonErrors.UnmarshalErrorCode))
		return
	}

	products, err := s.productService.SearchDefaults(productSearchRequest.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, commonErrors.NewError(commonErrors.UnknownErrorCode))
		return
	}

	c.JSON(http.StatusOK, products)
}
