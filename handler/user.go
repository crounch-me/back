package handler

import (
	"net/http"

	"github.com/Sehsyha/crounch-back/errorcode"
	"github.com/Sehsyha/crounch-back/model"

	"github.com/gin-gonic/gin"
)

func (hc *Context) Signup(c *gin.Context) {
	u := &model.User{}

	err := hc.UnmarshalPayload(c, u)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.UnmarshalCode, errorcode.UnmarshalDescription, http.StatusBadRequest)
		return
	}

	err = hc.Validate.Struct(u)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.InvalidCode, hc.GetValidationErrorDescription(err), http.StatusBadRequest)
		return
	}

	_, err = hc.Storage.GetUserByEmail(u.Email)
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
		hc.LogAndSendError(c, err, errorcode.DuplicateCode, errorcode.DuplicateDescription, http.StatusConflict)
		return
	}

	err = hc.Storage.CreateUser(u)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
		return
	}

	u.Password = nil

	c.JSON(http.StatusCreated, u)
}

func (hc *Context) Login(c *gin.Context) {
	u := &model.User{}

	err := hc.UnmarshalPayload(c, u)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.UnmarshalCode, errorcode.UnmarshalDescription, http.StatusBadRequest)
		return
	}

	err = hc.Validate.Struct(u)
	if err != nil {
		hc.LogAndSendError(c, err, errorcode.InvalidCode, hc.GetValidationErrorDescription(err), http.StatusBadRequest)
		return
	}

	authorization, err := hc.Storage.CreateAuthorization(u)

	if err != nil {
		if databaseError, ok := err.(*model.DatabaseError); ok {
			switch databaseError.Type {
			case model.ErrNotFound:
				c.AbortWithStatus(http.StatusForbidden)
				return
			case model.ErrWrongPassword:
				hc.LogAndSendError(c, err, errorcode.WrongPasswordCode, errorcode.WrongPasswordDescription, http.StatusBadRequest)
				return
			default:
				hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
				return
			}
		}
		hc.LogAndSendError(c, err, errorcode.DatabaseCode, errorcode.DatabaseDescription, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, authorization)
}
