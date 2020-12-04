package ports

import (
	"errors"
	"net/http"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/account/app"
	"github.com/crounch-me/back/internal/common/server"
	"github.com/crounch-me/back/util"
	"github.com/gin-gonic/gin"
)

const (
	userPath   = "/users"
	loginPath  = "/users/login"
	mePath     = "/me"
	logoutPath = "/logout"
)

type GinServer struct {
	accountService *app.AccountService
	validator      *util.Validator
}

func NewGinServer(accountService *app.AccountService, validator *util.Validator) (*GinServer, error) {
	if accountService == nil {
		return nil, errors.New("account gin server accountService is nil")
	}

	if validator == nil {
		return nil, errors.New("account gin server validator is nil")
	}

	return &GinServer{
		accountService: accountService,
		validator:      validator,
	}, nil
}

func (g *GinServer) ConfigureRoutes(r *gin.Engine) {
	r.POST(userPath, g.Signup)
	r.OPTIONS(userPath, server.OptionsHandler([]string{http.MethodPost}))

	r.POST(loginPath, g.Login)
	r.OPTIONS(loginPath, server.OptionsHandler([]string{http.MethodPost}))

	r.POST(logoutPath, g.Logout)
	r.OPTIONS(logoutPath, server.OptionsHandler([]string{http.MethodPost}))
}

func (s *GinServer) Signup(c *gin.Context) {
	signupRequest := &SignupRequest{}

	err := server.UnmarshalPayload(c.Request.Body, signupRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.UnmarshalErrorCode))
		return
	}

	err = s.validator.Struct(signupRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.InvalidErrorCode))
		return
	}

	err = s.accountService.Signup(signupRequest.Email, signupRequest.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, internal.NewError(internal.UnknownErrorCode))
		return
	}

	c.Status(http.StatusCreated)
}

func (s *GinServer) Login(c *gin.Context) {
	loginRequest := &LoginRequest{}

	err := server.UnmarshalPayload(c.Request.Body, loginRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.UnmarshalErrorCode))
		return
	}

	err = s.validator.Struct(loginRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.InvalidErrorCode))
		return
	}

	token, err := s.accountService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.UnknownErrorCode))
		return
	}

	tokenResponse := &TokenResponse{
		Token: token,
	}

	server.JSON(c, tokenResponse)
}

func (s *GinServer) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.Status(http.StatusNoContent)
		return
	}

	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, internal.NewError(internal.ForbiddenErrorCode))
		return
	}

	err = s.accountService.Logout(userUUID, token)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
