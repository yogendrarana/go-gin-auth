package routes

import (
	middleware "go-gin-auth/src/middlewares"

	"github.com/gin-gonic/gin"
)

func ProtectedRoutes(authGroup *gin.RouterGroup) {
	authGroup.GET("/test/protected", middleware.AuthMiddleware(), func(c *gin.Context) {
		user := c.MustGet("user")
		c.JSON(200, gin.H{
			"message": "Protected route test successful!",
			"user":    user,
		})
	})
}
