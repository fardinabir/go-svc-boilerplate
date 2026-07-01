package cases

import "github.com/labstack/echo/v4"

// RegisterRoutes registers case endpoints on the given group.
// The caller is responsible for scoping the group path (e.g. v1.Group("/cases")).
func RegisterRoutes(g *echo.Group, h Handler) {
	g.POST("", h.CreateCase)
	g.GET("", h.ListCases)
	g.GET("/:id", h.GetCaseByID)
	g.GET("/:id/assignee-email", h.GetAssigneeEmail)
}
