package ports

import (
	"errors"
	"net/http"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/common/utils"
	"github.com/crounch-me/back/internal/user/app"
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
	userService *app.UserService
	validator   *util.Validator
}

func NewGinServer(userService *app.UserService, validator *util.Validator) (*GinServer, error) {
	if userService == nil {
		return nil, errors.New("users gin server userService is nil")
	}

	if validator == nil {
		return nil, errors.New("users gin server validator is nil")
	}

	return &GinServer{
		userService: userService,
		validator:   validator,
	}, nil
}

func (g *GinServer) ConfigureRoutes(r *gin.Engine) {
	r.POST(userPath, g.Signup)
	r.OPTIONS(userPath, utils.OptionsHandler([]string{http.MethodPost}))

	r.POST(loginPath, g.Login)
	r.OPTIONS(loginPath, utils.OptionsHandler([]string{http.MethodPost}))
}

func (g *GinServer) Signup(c *gin.Context) {
	signupRequest := &SignupRequest{}

	err := utils.UnmarshalPayload(c.Request.Body, signupRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.UnmarshalErrorCode))
		return
	}

	err = g.validator.Struct(signupRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.InvalidErrorCode))
		return
	}

	err = g.userService.Signup(signupRequest.Email, signupRequest.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, internal.NewError(internal.UnknownErrorCode))
		return
	}

	c.Status(http.StatusCreated)
}

func (g *GinServer) Login(c *gin.Context) {
}
