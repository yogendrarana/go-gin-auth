// routes/authRoutes.go

package routes

import (
	"go-gin-auth/src/controllers"

	"github.com/gin-gonic/gin"
)

// InitializeAuthRoutes initializes authentication routes.
func InitializeAuthRoutes(authGroup *gin.RouterGroup) {
	{
		authGroup.POST("/auth/register", controllers.Register)
	}
}
