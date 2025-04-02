package middlewares

import (
	"go-gin-auth/src/models"
	"go-gin-auth/src/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthMiddleware is a middleware that checks if the user is authenticated
func AuthMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Authorization header is required."})
		return
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid authorization header format."})
		return
	}

	accessToken := authHeaderParts[1]
	isValidAccessToken, userIDPtr := services.ValidateJwtToken(accessToken, ctx)

	// if isValidAccessToken=false and userIDPtr=nil, then terminate the request
	if !isValidAccessToken && userIDPtr == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid access token."})
		return
	}

	// find the user in the database
	db := ctx.MustGet("db").(*gorm.DB)
	userID := *userIDPtr
	user, err := models.GetUserByID(db, userID)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to retrieve user"})
		return
	}

	if isValidAccessToken {
		ctx.Set("user", user)
		ctx.Next()
	}
}
