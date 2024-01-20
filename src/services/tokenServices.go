package services

import (
	"fmt"
	"go-gin-auth/src/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 10
)

// generate access token
func GenerateAccessToken(userID uint) (string, error) {
	// access token expires in 15 minutes
	accessExpiry := time.Now().Add(time.Minute * 15).Unix()

	// defining the claims for the access tokens
	accessClaims := jwt.MapClaims{"exp": accessExpiry, "sub": userID}

	// create the access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSignedToken, err := accessToken.SignedString([]byte(os.Getenv("ACCESS_JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return accessSignedToken, nil
}

// check the validity of the access token
func ValidateJwtToken(accessToken string, ctx *gin.Context) (bool, *uint) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// not required to check signing method but checking it is a good practice
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_JWT_SECRET")), nil
	})

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract claims"})
		return false, nil
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract user ID from claims"})
		return false, nil
	}

	userIDInt := uint(userIDFloat)

	if err != nil && !token.Valid {
		fmt.Println("Error parsing JWT:", err)
		// cannot abort the request here because the refresh token still might be valid
		// ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		return false, &userIDInt
	}

	return true, &userIDInt
}

// generate refresh token
func GenerateRefreshTokenAndHash() (string, string, error) {
	refreshToken := uuid.New().String()

	// Hash the refresh token using bcrypt
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcryptCost)
	if err != nil {
		return "", "", err
	}

	return refreshToken, string(hashedToken), nil
}

// check validity of refresh token
func ValidateRefreshToken(refreshToken string, hashedTokens []models.RefreshToken) (bool, error) {
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
