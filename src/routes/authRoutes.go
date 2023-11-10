package routes

import (
	"go-gin-auth/src/handlers"
	middleware "go-gin-auth/src/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeAuthRoutes(authGroup *gin.RouterGroup) {
	// middleware
	authGroup.Use(middleware.DatabaseMiddleware())

	// routes
	authGroup.POST("/auth/register", handlers.Register)
	authGroup.POST("/auth/login", handlers.Login)
	authGroup.POST("/auth/delete", handlers.DeleteUser)
}
