package routes

import (
	"go-gin-auth/src/handlers"
	middlewares "go-gin-auth/src/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(authGroup *gin.RouterGroup) {
	// middleware
	authGroup.Use(middlewares.DatabaseMiddleware())

	// routes
	authGroup.POST("/auth/register", handlers.Register)
	authGroup.POST("/auth/login", handlers.Login)
	authGroup.POST("/auth/delete", handlers.DeleteUser)
}
