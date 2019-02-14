package handler

import (
	"net/http"

	"github.com/Sehsyha/crounch-back/errorcode"
	"github.com/Sehsyha/crounch-back/model"

	"github.com/gin-gonic/gin"
)

func (hc *Context) Signup(c *gin.Context) {
	u := &model.User{}

	if hc.UnmarshalPayload(c, u) {
		return
	}

	_, err := hc.Storage.GetUserByEmail(u.Email)
	if err != nil {
		if databaseError, ok := err.(*model.DatabaseError); ok {
			switch databaseError.Type {
			case model.ErrNotFound:
				break
			default:
				hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
				return
			}
		} else {
			hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
			return
		}
	}

	if err == nil {
		hc.LogAndSendError(c, err, errorcode.Duplicate, "User with this email already exists", http.StatusConflict)
		return
	}

	err = hc.Storage.CreateUser(u)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, u)
}

func (hc *Context) Login(c *gin.Context) {
	u := &model.User{}

	if hc.UnmarshalPayload(c, u) {
		return
	}

	authorization, err := hc.Storage.CreateAuthorization(u)

	if err != nil {
		hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, authorization)
}
