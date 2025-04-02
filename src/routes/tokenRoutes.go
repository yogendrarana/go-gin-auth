package routes

import (
	"go-gin-auth/src/handlers"

	"github.com/gin-gonic/gin"
)

func TokenRoutes(tokenGroup *gin.RouterGroup) {
	tokenGroup.GET("/token/new-access-token", handlers.HandleRefreshToken)
}
