package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/util"
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
		return "", internal.NewError(internal.UnknownErrorCode).WithCall("utils", "GetUserIDFromContext")
	}

	return userID.(string), nil
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
	allowedHeaders := []string{util.HeaderContentType, util.HeaderAuthorization, util.HeaderAccept}
	return func(c *gin.Context) {
		c.Writer.Header().Set(util.HeaderAccessControlAllowOrigin, "*")
		c.Writer.Header().Set(util.HeaderAccessControlAllowMethods, strings.Join(allowedMethods, ","))
		c.Writer.Header().Set(util.HeaderAccessControlAllowHeaders, strings.Join(allowedHeaders, ","))
		c.AbortWithStatus(http.StatusOK)
	}
}
