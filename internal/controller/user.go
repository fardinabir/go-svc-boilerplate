package controller

import (
	"net/http"
	"strconv"

	"github.com/fardinabir/go-svc-boilerplate/internal/errors"
	"github.com/fardinabir/go-svc-boilerplate/internal/model"
	"github.com/fardinabir/go-svc-boilerplate/internal/service"
	"github.com/labstack/echo/v4"
)

// UserHandler is the request handler for the user endpoint.
type UserHandler interface {
	CreateUser(c echo.Context) error
	ListUsers(c echo.Context) error
	GetUserByID(c echo.Context) error
}

type userHandler struct {
	Handler
	service service.UserService
}

// NewUserHandler returns a new instance of the user handler.
func NewUserHandler(s service.UserService) UserHandler {
	return &userHandler{service: s}
}

// UserCreateRequest represents the request for creating a user
type UserCreateRequest struct {
	Name  string `json:"name" validate:"required,validUserName"`
	Email string `json:"email" validate:"required,email"`
}

// @Summary    Create a user
// @Tags       users
// @Accept     json
// @Produce    json
// @Param      request body        UserCreateRequest  true  "User create request"
// @Success    201     {object}    ResponseData{data=model.User}
// @Failure    400     {object}    ResponseError
// @Failure    500     {object}    ResponseError
// @Router     /users [post]
func (h *userHandler) CreateUser(c echo.Context) error {
	var req UserCreateRequest
	if err := h.MustBind(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: err.Error()}}})
	}

	u := &model.User{Name: req.Name, Email: req.Email}
	if err := h.service.CreateUser(u); err != nil {
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}
	return c.JSON(http.StatusCreated, ResponseData{Data: u})
}

// @Summary    List users
// @Tags       users
// @Produce    json
// @Success    200     {object}    ResponseData{data=[]model.User}
// @Failure    500     {object}    ResponseError
// @Router     /users [get]
func (h *userHandler) ListUsers(c echo.Context) error {
	users, err := h.service.ListUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}
	return c.JSON(http.StatusOK, ResponseData{Data: users})
}

// @Summary    Get user by ID
// @Tags       users
// @Produce    json
// @Param      id     path        int  true  "User ID"
// @Success    200     {object}    ResponseData{data=model.User}
// @Failure    400     {object}    ResponseError
// @Failure    404     {object}    ResponseError
// @Failure    500     {object}    ResponseError
// @Router     /users/{id} [get]
func (h *userHandler) GetUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest,
			ResponseError{Errors: []Error{{Code: errors.CodeBadRequest, Message: "invalid id"}}})
	}
	user, err := h.service.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ResponseError{Errors: []Error{{Code: errors.CodeInternalServerError, Message: err.Error()}}})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound,
			ResponseError{Errors: []Error{{Code: errors.CodeNotFound, Message: "user not found"}}})
	}
	return c.JSON(http.StatusOK, ResponseData{Data: user})
}
