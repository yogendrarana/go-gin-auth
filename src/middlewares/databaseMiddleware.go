package middlewares

import (
	"go-gin-auth/src/initializers"

	"github.com/gin-gonic/gin"
)

func DatabaseMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := initializers.GetDB()
		ctx.Set("db", db)
		ctx.Next()
	}
}
