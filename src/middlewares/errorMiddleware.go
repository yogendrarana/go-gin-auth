package middlewares

import (
	custom_errors "go-gin-auth/src/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// error handler middleware
func ErrorMiddleware(c *gin.Context) {
	c.Next()
	for _, err := range c.Errors {
		switch e := err.Err.(type) {
		case *custom_errors.AppError:
			c.AbortWithStatusJSON(e.Code, gin.H{"success": false, "message": e.Message})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Internal server error"})
		}
	}
}
