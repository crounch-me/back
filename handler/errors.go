package handler

import (
	"net/http"

	"github.com/crounch-me/back/internal/account"
	"github.com/crounch-me/back/internal/common/errors"
	list "github.com/crounch-me/back/internal/listing"
	"github.com/crounch-me/back/internal/products"
)

var errorToStatus map[string]int

func initializeErrorsMap() {
	errorToStatus = make(map[string]int)

	// Internal errors
	errorToStatus[errors.UnmarshalErrorCode] = http.StatusBadRequest
	errorToStatus[errors.InvalidErrorCode] = http.StatusBadRequest
	errorToStatus[errors.ForbiddenErrorCode] = http.StatusForbidden
	errorToStatus[errors.UnauthorizedErrorCode] = http.StatusUnauthorized
	errorToStatus[errors.UnknownErrorCode] = http.StatusInternalServerError

	// User errors
	errorToStatus[account.UserNotFoundErrorCode] = http.StatusNotFound
	errorToStatus[account.DuplicateUserErrorCode] = http.StatusConflict

	// List errors
	errorToStatus[list.ListNotFoundErrorCode] = http.StatusNotFound
	errorToStatus[list.DuplicateProductInListErrorCode] = http.StatusConflict

	// Product in list errors
	errorToStatus[list.ProductInListNotFoundErrorCode] = http.StatusNotFound

	// Product errors
	errorToStatus[products.ProductNotFoundErrorCode] = http.StatusNotFound
}

// ErrorCodeToHTTPStatus returns the status from the given error
func (hc *Context) ErrorCodeToHTTPStatus(code string) int {
	if errorToStatus == nil {
		initializeErrorsMap()
	}

	status, ok := errorToStatus[code]
	if !ok {
		return http.StatusInternalServerError
	}

	return status
}
