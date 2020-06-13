package handler

import (
	"net/http"

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
