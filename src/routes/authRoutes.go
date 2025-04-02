package routes

import (
	"go-gin-auth/src/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(authGroup *gin.RouterGroup) {
	authGroup.POST("/auth/register", handlers.Register)
	authGroup.POST("/auth/login", handlers.Login)
	authGroup.POST("/auth/logout", handlers.Logout)
}
