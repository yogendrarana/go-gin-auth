package middleware

import (
	"errors"
	"fmt"
	"go-gin-auth/src/models"
	"go-gin-auth/src/services"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		accessToken := authHeaderParts[1]
		isValidAccess, user := isValidAccessToken(accessToken, c)

		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})

		if isValidAccess && user != nil {
			c.Set("user", user)
			c.Next()
		}

		if !isValidAccess && user == nil {
			isValidRefresh := isValidRefreshToken(c)

			if isValidRefresh {
				// Create a new access token
				newAccessToken, err := services.GenerateAccessToken(user.ID)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
					return
				}

				// Set the new access token in the response header
				c.Header("Authorization", fmt.Sprintf("Bearer %s", newAccessToken))
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token. Please login again."})
				return
			}
		}

	}
}

// check validity of access token
func isValidAccessToken(accessToken string, c *gin.Context) (bool, *models.User) {
	// get db connection
	db := c.MustGet("db").(*gorm.DB)

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// not required to check signing method but a good practice
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_JWT_SECRET")), nil
	})

	if token == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		return false, nil
	}

	// If there is any error other than TokenExpired, return 401 Unauthorized
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		return false, nil
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to access claims"})
		return false, nil
	}

	userID, ok := claims["sub"].(uint)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract user ID from claims"})
		return false, nil
	}

	// query user from database and check if user
	var user models.User
	result := db.Where("id = ?", userID).First(&user)

	if result.Error != nil || result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return false, nil
	}

	return true, &user
}

// check validity if refresh token
func isValidRefreshToken(c *gin.Context) bool {
	db := c.MustGet("db").(*gorm.DB)

	// get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return false
	}

	// get refresh token from database
	var refreshTokenFromDB models.RefreshToken
	result := db.Where("token_hash = ?", refreshToken).First(&refreshTokenFromDB)

	if result.Error != nil || result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return false
	}

	// check if refresh token is expired
	if refreshTokenFromDB.ExpiresAt < (time.Now().Unix() * 1000) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return false
	}

	return true
}
