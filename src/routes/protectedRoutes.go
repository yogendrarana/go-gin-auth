package routes

import (
	middleware "go-gin-auth/src/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeProtectedRoutes(authGroup *gin.RouterGroup) {
	// test routes
	authGroup.GET("/test/protected", middleware.AuthMiddleware(), func(c *gin.Context) {
		// get the user from the context
		user := c.MustGet("user")
		c.JSON(200, gin.H{
			"message": "Protected route test successful!",
			"user":    user,
		})
	})
}
