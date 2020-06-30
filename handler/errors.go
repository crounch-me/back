package handler

import (
	"net/http"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/authorization.go"
	"github.com/crounch-me/back/domain/lists"
	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
)

var errorToStatus map[string]int

func initializeErrorsMap() {
	errorToStatus = make(map[string]int)

	// Domain errors
	errorToStatus[domain.UnmarshalErrorCode] = http.StatusBadRequest
	errorToStatus[domain.InvalidErrorCode] = http.StatusBadRequest
	errorToStatus[domain.ForbiddenErrorCode] = http.StatusForbidden
	errorToStatus[domain.UnauthorizedErrorCode] = http.StatusUnauthorized
	errorToStatus[domain.UnknownErrorCode] = http.StatusInternalServerError

	// User errors
	errorToStatus[users.UserNotFoundErrorCode] = http.StatusNotFound
	errorToStatus[users.DuplicateUserErrorCode] = http.StatusConflict

	// List errors
	errorToStatus[lists.ListNotFoundErrorCode] = http.StatusNotFound
	errorToStatus[lists.DuplicateProductInListErrorCode] = http.StatusConflict

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
