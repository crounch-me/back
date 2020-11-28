package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

const (
	HeaderContentLocation = "Content-Location"
)

func GetUserIDFromContext(c *gin.Context) (string, error) {
	return "", nil
}

func UnmarshalPayload(payload io.ReadCloser, i interface{}) error {
	bytePayload, err := ioutil.ReadAll(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytePayload, i)
}
