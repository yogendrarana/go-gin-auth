package routes

import (
	"go-gin-auth/src/handlers"
	"go-gin-auth/src/middlewares"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(authGroup *gin.RouterGroup) {
	authGroup.GET("/admin", middlewares.AuthMiddleware, handlers.GetDashboardData)
}
