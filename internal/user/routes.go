package user

import "github.com/labstack/echo/v4"

// RegisterRoutes registers user endpoints on the given group.
// The caller is responsible for scoping the group path (e.g. v1.Group("/users")).
func RegisterRoutes(g *echo.Group, h Handler) {
	g.POST("", h.CreateUser)
	g.GET("", h.ListUsers)
	g.GET("/:id", h.GetUserByID)
}
