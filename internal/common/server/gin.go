package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/crounch-me/back/internal/account"
	accountApp "github.com/crounch-me/back/internal/account/app"
	"github.com/crounch-me/back/internal/common/errors"
	"github.com/gin-gonic/gin"
)

const (
	HeaderContentLocation = "Content-Location"

	ContextUserID = "ContextUserID"
)

type DataResponse struct {
	Data interface{} `json:"data"`
}

func NewDataResponse(data interface{}) *DataResponse {
	return &DataResponse{
		Data: data,
	}
}

func GetUserIDFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get(ContextUserID)
	if !exists {
		return "", errors.NewError(errors.UnknownErrorCode).WithCall("utils", "GetUserIDFromContext")
	}

	return userID.(string), nil
}

func CheckUserAuthorization(accountService *accountApp.AccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			fmt.Println("Unauthorized - No token provided")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userUUID, err := accountService.GetUserUUIDByToken(token)
		if err != nil {
			if err.Error() == account.UserNotFoundErrorCode {
				c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewError(errors.UnauthorizedErrorCode))
				return
			}

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Set(ContextUserID, userUUID)
	}
}

func UnmarshalPayload(payload io.ReadCloser, i interface{}) error {
	bytePayload, err := ioutil.ReadAll(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytePayload, i)
}

func JSON(c *gin.Context, response interface{}) {
	dataResponse := NewDataResponse(response)
	c.JSON(http.StatusOK, dataResponse)
}

func OptionsHandler(allowedMethods []string) gin.HandlerFunc {
	allowedMethods = append(allowedMethods, http.MethodOptions)
	allowedHeaders := []string{HeaderContentType, HeaderAuthorization, HeaderAccept}

	return func(c *gin.Context) {
		c.Writer.Header().Set(HeaderAccessControlAllowOrigin, "*")
		c.Writer.Header().Set(HeaderAccessControlAllowMethods, strings.Join(allowedMethods, ","))
		c.Writer.Header().Set(HeaderAccessControlAllowHeaders, strings.Join(allowedHeaders, ","))
		c.AbortWithStatus(http.StatusOK)
	}
}
