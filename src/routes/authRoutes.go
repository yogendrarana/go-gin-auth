// routes/authRoutes.go

package routes

import (
	"go-gin-auth/src/controllers"
	middleware "go-gin-auth/src/middlewares"

	"github.com/gin-gonic/gin"
)

// InitializeAuthRoutes initializes authentication routes.
func InitializeAuthRoutes(authGroup *gin.RouterGroup) {
	// middleware
	authGroup.Use(middleware.DatabaseMiddleware())

	// routes
	authGroup.POST("/auth/register", controllers.Register)
	authGroup.POST("/auth/login", controllers.Login)

}
