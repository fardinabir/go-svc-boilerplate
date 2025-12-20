package controller

import "github.com/labstack/echo/v4"

func InitRoutes(api *echo.Group, controller UserHandler) {
    users := api.Group("/users")
    {
        users.POST("", controller.CreateUser)
        users.GET("", controller.ListUsers)
        users.GET("/:id", controller.GetUserByID)
    }
}
