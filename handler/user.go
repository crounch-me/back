package handler

import (
	"net/http"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/user"

	"github.com/gin-gonic/gin"
)

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gt=3"`
}

// Signup creates a new user
// @Summary Creates a new user
// @ID signup
// @Tags user
// @Accept json
// @Produce  json
// @Param user body SignupRequest true "User to signup with"
// @Success 200 {object} user.User
// @Failure 500 {object} internal.Error
// @Router /users/signup [post]
func (hc *Context) Signup(c *gin.Context) {
	u := &user.User{}

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

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gt=3"`
}

// Login creates a new user authorization if is found and password is good
// @Summary Creates a new user authorization
// @ID login
// @Tags user
// @Accept json
// @Produce  json
// @Param user body LoginRequest true "User to login with"
// @Success 200 {object} authorization.Authorization
// @Failure 500 {object} internal.Error
// @Router /users/login [post]
func (hc *Context) Login(c *gin.Context) {
	u := &LoginRequest{}

	err := hc.UnmarshalAndValidate(c, u)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	authorization, err := hc.Services.Authorization.CreateAuthorization(u.Email, u.Password)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, authorization)
}

// Logout removes the user authorization if it is found
// @Summary Removes an user authorization
// @ID logout
// @Tags user
// @Success 204
// @Failure 500 {object} internal.Error
// @Security ApiKeyAuth
// @Router /logout [post]
func (hc *Context) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		hc.LogAndSendError(c, internal.NewError(internal.UnauthorizedErrorCode))
		return
	}

	err := hc.Services.Authorization.Logout(token)
	if err != nil {
		hc.LogAndSendError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Me returns user informations
// @Summary Removes an user authorization
// @ID me
// @Tags user
// @Produce json
// @Success 200 {object} users.User
// @Failure 500 {object} internal.Error
// @Security ApiKeyAuth
// @Router /me [get]
func (hc *Context) Me(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		hc.LogAndSendError(c, internal.NewError(internal.UnauthorizedErrorCode))
		return
	}

	u, err := hc.Services.User.GetByToken(token)

	if err != nil {
		if err.Code == user.UserNotFoundErrorCode {
			hc.LogAndSendError(c, internal.NewError(internal.UnauthorizedErrorCode))
			return
		}
		hc.LogAndSendError(c, err)
		return
	}

	u.Password = nil

	c.JSON(http.StatusOK, u)
}
