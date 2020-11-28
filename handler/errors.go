package handler

import (
	"net/http"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/authorization.go"
	"github.com/crounch-me/back/internal/list"
	"github.com/crounch-me/back/internal/products"
	"github.com/crounch-me/back/internal/users"
)

var errorToStatus map[string]int

func initializeErrorsMap() {
	errorToStatus = make(map[string]int)

	// Internal errors
	errorToStatus[internal.UnmarshalErrorCode] = http.StatusBadRequest
	errorToStatus[internal.InvalidErrorCode] = http.StatusBadRequest
	errorToStatus[internal.ForbiddenErrorCode] = http.StatusForbidden
	errorToStatus[internal.UnauthorizedErrorCode] = http.StatusUnauthorized
	errorToStatus[internal.UnknownErrorCode] = http.StatusInternalServerError

	// User errors
	errorToStatus[users.UserNotFoundErrorCode] = http.StatusNotFound
	errorToStatus[users.DuplicateUserErrorCode] = http.StatusConflict

	// List errors
	errorToStatus[list.ListNotFoundErrorCode] = http.StatusNotFound
	errorToStatus[list.DuplicateProductInListErrorCode] = http.StatusConflict

	// Product in list errors
	errorToStatus[list.ProductInListNotFoundErrorCode] = http.StatusNotFound

	// Product errors
	errorToStatus[products.ProductNotFoundErrorCode] = http.StatusNotFound

	// Authorization errors
	errorToStatus[authorization.WrongPasswordErrorCode] = http.StatusBadRequest
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
