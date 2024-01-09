package middlewares

import (
	"fmt"
	"go-gin-auth/src/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		accessToken := authHeaderParts[1]
		isValidAccessToken, userIDPtr := validateAccessToken(accessToken, ctx)

		// if isValidAccessToken=false and userIDPtr=nil, then terminate the request
		if !isValidAccessToken && userIDPtr == nil {
			return
		}

		// find the user in the database
		db := ctx.MustGet("db").(*gorm.DB)
		userID := *userIDPtr
		user, err := models.GetUserByID(db, userID)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		if isValidAccessToken {
			ctx.Set("user", user)
			ctx.Next()
		}
	}
}

// check validity of access token
func validateAccessToken(accessToken string, c *gin.Context) (bool, *uint) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// not required to check signing method but a good practice
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_JWT_SECRET")), nil
	})

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract claims"})
		return false, nil
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract user ID from claims"})
		return false, nil
	}

	userIDInt := uint(userIDFloat)

	if err != nil && !token.Valid {
		fmt.Println("Error parsing JWT:", err)
		// cannot abort the request here because the refresh token still might be valid
		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		return false, &userIDInt
	}

	return true, &userIDInt
}

// check validity of refresh token
func validateRefreshToken(refreshToken string, hashedTokens []models.RefreshToken) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedTokens[0].TokenHash), []byte(refreshToken))
	if err != nil {
		return false, err
	}

	// Check if the refresh token has expired
	if hashedTokens[0].ExpiresAt < time.Now().Unix() {
		return false, nil
	}

	return true, nil
}
