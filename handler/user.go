package handler

import (
	"net/http"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"

	"github.com/gin-gonic/gin"
)

func (hc *Context) Signup(c *gin.Context) {
	u := &users.User{}

	err := hc.UnmarshalAndValidate(c, u)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	user, err := hc.Services.User.CreateUser(u.Email, *u.Password)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	user.Password = nil

	c.JSON(http.StatusCreated, user)
}

func (hc *Context) Login(c *gin.Context) {
	u := &users.User{}

	err := hc.UnmarshalAndValidate(c, u)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	authorization, err := hc.Services.Authorization.CreateAuthorization(u.Email, *u.Password)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, authorization)
}

func (hc *Context) Me(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		hc.LogAndSendError(c, domain.NewError(domain.UnauthorizedErrorCode))
		return
	}

	user, err := hc.Services.User.GetByToken(token)

	if err != nil {
		if err.Code == users.UserNotFoundErrorCode {
			hc.LogAndSendError(c, domain.NewError(domain.UnauthorizedErrorCode))
			return
		}
		hc.LogAndSendError(c, err)
		return
	}

	user.Password = nil

	c.JSON(http.StatusOK, user)
}
