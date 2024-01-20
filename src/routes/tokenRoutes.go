package routes

import (
	"go-gin-auth/src/handlers"

	"github.com/gin-gonic/gin"
)

func TokenRoutes(tokenGroup *gin.RouterGroup) {
	tokenGroup.GET("/token/refresh", handlers.HandleRefreshToken)
}
