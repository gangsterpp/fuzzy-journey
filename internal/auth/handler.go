package auth

import (
	"errors"
	"net/http"

	"github.com/gangsterpp/fuzzy-journey/internal/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	Delete(c *gin.Context)
}

type Handler struct {
	service AuthService
}

func (h *Handler) Login(c *gin.Context) {

	var auth AuthModel
	if err := c.ShouldBind(&auth); err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			response.CodeInvalidRequest,
			err.Error(),
		)
		return
	}

	user, err := h.service.Login(c.Request.Context(), &auth)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrInvalidCredentials):
			response.Fail(c, http.StatusUnauthorized, response.CodeInvalidCredentials, err.Error())
		default:
			response.Fail(
				c,
				http.StatusInternalServerError,
				response.CodeInternal,
				err.Error(),
			)
		}
		return
	}
	response.OK(c, user)

}

func (h *Handler) Register(c *gin.Context) {
	var auth AuthModel
	if err := c.ShouldBind(&auth); err != nil {
		response.Fail(
			c,
			http.StatusBadRequest,
			response.CodeInternal,
			err.Error(),
		)
		return
	}
	user, err := h.service.Register(c.Request.Context(), &auth)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrUserAlreadyExists):
			response.Fail(c, http.StatusNotFound, response.CodeUserAlreadyExists, err.Error())
		default:
			response.Fail(
				c,
				http.StatusInternalServerError,
				response.CodeInternal,
				err.Error(),
			)
		}
		return
	}

	response.OK(c, user)

}

func (h *Handler) Delete(c *gin.Context) {
	value, exist := c.GetQuery("id")
	if !exist {
		return
	}
	user, err := h.service.Delete(c.Request.Context(), value)

	if err != nil {
		switch {
		case errors.Is(err, response.ErrUserNotFound):
			response.Fail(c, http.StatusNotFound, response.CodeErrUserNotFound, err.Error())
		default:
			response.Fail(
				c,
				http.StatusInternalServerError,
				response.CodeInternal,
				err.Error(),
			)
		}
		return
	}
	response.OK(c, user)

}

var _ AuthHandler = (*Handler)(nil)

func CreateAuthHandler(service AuthService) AuthHandler {

	return &Handler{service: service}
}
