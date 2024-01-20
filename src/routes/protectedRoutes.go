package routes

import (
	middleware "go-gin-auth/src/middlewares"

	"github.com/gin-gonic/gin"
)

func ProtectedRoutes(authGroup *gin.RouterGroup) {
	authGroup.GET("/protected", middleware.AuthMiddleware(), func(ctx *gin.Context) {
		user := ctx.MustGet("user")
		ctx.JSON(200, gin.H{
			"message": "Protected route test successful!",
			"user":    user,
		})
	})
}
