package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/crounch-me/back/internal"
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
