// Package health provides the service health endpoint.
package health

import (
	"net/http"
	"time"

	"github.com/fardinabir/go-svc-boilerplate/pkg/response"
	"github.com/labstack/echo/v4"
)

// Handler is the request handler for the health endpoint.
type Handler interface {
	Health(c echo.Context) error
}

type handler struct{}

// New returns a new instance of the health handler.
func New() Handler {
	return &handler{}
}

// RegisterRoutes mounts the health endpoint under the given group.
func RegisterRoutes(g *echo.Group, h Handler) {
	g.GET("/health", h.Health)
}

// @Summary	Health check
// @Tags		health
// @Produce	json
// @Success	200	{object}	response.ResponseData{data=time.Time}
// @Router		/health [get]
func (h *handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, response.ResponseData{
		Data: map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now(),
		},
	})
}
