package middlewares

import (
	"fmt"
	"go-gin-auth/src/models"
	"go-gin-auth/src/services"
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
		isValidAccessToken, userIDPtr := checkAccessTokenValidity(accessToken, c)

		// if isValidAccessToken=false and userIDPtr=nil, then terminate the request
		if !isValidAccessToken && userIDPtr == nil {
			return
		}

		// find the user in the database
		db := c.MustGet("db").(*gorm.DB)
		userID := *userIDPtr
		user, err := models.GetUserByID(db, userID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		if isValidAccessToken {
			c.Set("user", user)
			c.Next()
		}

		if !isValidAccessToken {
			refreshToken, err := c.Cookie("refresh_token")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
				return
			}

			// check validity of refresh token
			isValidRefreshToken, err := checkRefreshTokenValidity(refreshToken, user.RefreshTokens)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to compare refresh token hashes"})
				return
			}

			// if valid refresh token issue new access token
			if isValidRefreshToken {
				// generate access tokens
				accessSignedToken, err := services.GenerateAccessToken(user.ID)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"new_access_token": accessSignedToken})

				// set the user in the Gin context
				c.Set("user", user)

				// Call c.Next() to pass the request to the next handler
				c.Next()
			}

			// if invalid refresh token terminate the request
			if !isValidRefreshToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
				return
			}
		}
	}
}

// check validity of access token
func checkAccessTokenValidity(accessToken string, c *gin.Context) (bool, *uint) {
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
func checkRefreshTokenValidity(refreshToken string, hashedTokens []models.RefreshToken) (bool, error) {
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
