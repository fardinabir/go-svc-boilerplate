package user

import (
	"net/http"
	"strconv"

	apierr "github.com/fardinabir/go-svc-boilerplate/internal/errors"
	"github.com/fardinabir/go-svc-boilerplate/pkg/response"
	"github.com/fardinabir/go-svc-boilerplate/pkg/web"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Handler is the request handler for the user endpoints.
type Handler interface {
	CreateUser(c echo.Context) error
	ListUsers(c echo.Context) error
	GetUserByID(c echo.Context) error
}

type handler struct {
	web.Base
	service Service
}

// NewHandler returns a new instance of the user handler.
func NewHandler(s Service) Handler {
	return &handler{service: s}
}

// CreateRequest represents the request for creating a user.
type CreateRequest struct {
	Name  string `json:"name" validate:"required,validUserName"`
	Email string `json:"email" validate:"required,email"`
}

// @Summary    Create a user
// @Tags       users
// @Accept     json
// @Produce    json
// @Param      request body        user.CreateRequest  true  "User create request"
// @Success    201     {object}    response.ResponseData{data=user.User}
// @Failure    400     {object}    response.APIError
// @Failure    500     {object}    response.APIError
// @Router     /users [post]
func (h *handler) CreateUser(c echo.Context) error {
	var req CreateRequest
	if err := h.MustBind(c, &req); err != nil {
		return response.Respond(c, apierr.ErrBadRequest, err.Error())
	}
	u := &User{Name: req.Name, Email: req.Email}
	if err := h.service.CreateUser(u); err != nil {
		return response.Respond(c, apierr.ErrInternalServerError)
	}
	return c.JSON(http.StatusCreated, response.ResponseData{Data: u})
}

// @Summary    List users
// @Tags       users
// @Produce    json
// @Success    200     {object}    response.ResponseData{data=[]user.User}
// @Failure    500     {object}    response.APIError
// @Router     /users [get]
func (h *handler) ListUsers(c echo.Context) error {
	users, err := h.service.ListUsers()
	if err != nil {
		return response.Respond(c, apierr.ErrInternalServerError)
	}
	return c.JSON(http.StatusOK, response.ResponseData{Data: users})
}

// @Summary    Get user by ID
// @Tags       users
// @Produce    json
// @Param      id     path        int  true  "User ID"
// @Success    200     {object}    response.ResponseData{data=user.User}
// @Failure    400     {object}    response.APIError
// @Failure    404     {object}    response.APIError
// @Failure    500     {object}    response.APIError
// @Router     /users/{id} [get]
func (h *handler) GetUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return response.Respond(c, apierr.ErrBadRequest, "invalid id")
	}
	u, err := h.service.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Respond(c, ErrNotFound)
		}
		return response.Respond(c, apierr.ErrInternalServerError)
	}
	return c.JSON(http.StatusOK, response.ResponseData{Data: u})
}
