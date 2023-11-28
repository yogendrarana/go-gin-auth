package middlewares

import (
	"go-gin-auth/src/initializers"

	"github.com/gin-gonic/gin"
)

func DatabaseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := initializers.GetDB()
		c.Set("db", db)
		c.Next()
	}
}
