package cases

import (
	"net/http"
	"strconv"

	apierr "github.com/fardinabir/go-svc-boilerplate/internal/errors"
	"github.com/fardinabir/go-svc-boilerplate/pkg/response"
	"github.com/fardinabir/go-svc-boilerplate/pkg/web"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Handler is the request handler for the case endpoints.
type Handler interface {
	CreateCase(c echo.Context) error
	ListCases(c echo.Context) error
	GetCaseByID(c echo.Context) error
	GetAssigneeEmail(c echo.Context) error
}

type handler struct {
	web.Base
	service Service
}

// NewHandler returns a new instance of the cases handler.
func NewHandler(s Service) Handler {
	return &handler{service: s}
}

// CreateRequest represents the request for creating a case.
type CreateRequest struct {
	FileNumber string `json:"file_number" validate:"required"`
	Status     string `json:"status" validate:"required"`
	ServicerID int    `json:"servicer_id" validate:"required"`
	AssigneeID int    `json:"assignee_id" validate:"required"`
}

// @Summary    Create a case
// @Tags       cases
// @Accept     json
// @Produce    json
// @Param      request body        cases.CreateRequest  true  "Case create request"
// @Success    201     {object}    response.ResponseData{data=cases.Case}
// @Failure    400     {object}    response.APIError
// @Failure    500     {object}    response.APIError
// @Router     /cases [post]
func (h *handler) CreateCase(c echo.Context) error {
	var req CreateRequest
	if err := h.MustBind(c, &req); err != nil {
		return response.Respond(c, apierr.ErrBadRequest, err.Error())
	}
	cs := &Case{FileNumber: req.FileNumber, Status: req.Status, ServicerID: req.ServicerID, AssigneeID: req.AssigneeID}
	if err := h.service.CreateCase(cs); err != nil {
		return response.Respond(c, apierr.ErrInternalServerError)
	}
	return c.JSON(http.StatusCreated, response.ResponseData{Data: cs})
}

// @Summary    List cases
// @Tags       cases
// @Produce    json
// @Success    200     {object}    response.ResponseData{data=[]cases.Case}
// @Failure    500     {object}    response.APIError
// @Router     /cases [get]
func (h *handler) ListCases(c echo.Context) error {
	cs, err := h.service.ListCases()
	if err != nil {
		return response.Respond(c, apierr.ErrInternalServerError)
	}
	return c.JSON(http.StatusOK, response.ResponseData{Data: cs})
}

// @Summary    Get case by ID
// @Tags       cases
// @Produce    json
// @Param      id     path        int  true  "Case ID"
// @Success    200     {object}    response.ResponseData{data=cases.Case}
// @Failure    400     {object}    response.APIError
// @Failure    404     {object}    response.APIError
// @Failure    500     {object}    response.APIError
// @Router     /cases/{id} [get]
func (h *handler) GetCaseByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return response.Respond(c, apierr.ErrBadRequest, "invalid id")
	}
	cs, err := h.service.GetCaseByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Respond(c, ErrNotFound)
		}
		return response.Respond(c, apierr.ErrInternalServerError)
	}
	return c.JSON(http.StatusOK, response.ResponseData{Data: cs})
}

// @Summary    Get the email of the case assignee (cross-domain read)
// @Tags       cases
// @Produce    json
// @Param      id     path        int  true  "Case ID"
// @Success    200     {object}    response.ResponseData{data=string}
// @Failure    400     {object}    response.APIError
// @Failure    404     {object}    response.APIError
// @Failure    500     {object}    response.APIError
// @Router     /cases/{id}/assignee-email [get]
func (h *handler) GetAssigneeEmail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return response.Respond(c, apierr.ErrBadRequest, "invalid id")
	}
	email, err := h.service.AssigneeEmail(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Respond(c, ErrNotFound)
		}
		return response.Respond(c, apierr.ErrInternalServerError)
	}
	return c.JSON(http.StatusOK, response.ResponseData{Data: email})
}
