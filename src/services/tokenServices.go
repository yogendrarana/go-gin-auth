package services

import (
	"errors"
	"fmt"
	"go-gin-auth/src/models"
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

	if err != nil {
		return false, nil
	}
	fmt.Println("tokennnnnn", token)

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return false, nil
	}

	userIDInt := uint(userIDFloat)

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

func ValidateRefreshToken(refreshToken string, refreshTokens []models.RefreshToken) (*models.RefreshToken, error) {
	var matchingToken *models.RefreshToken

	// Loop through the array of refresh tokens
	for _, tkn := range refreshTokens {
		// Compare the token in the cookie with the hashed token in the database
		err := bcrypt.CompareHashAndPassword([]byte(tkn.TokenHash), []byte(refreshToken))
		if err != nil {
			fmt.Println("error hai", err)
		}

		if err == nil {
			// If the comparison is successful, store the matching refresh token
			matchingToken = &tkn
			break
		}
	}

	// Check if a matching token was found
	if matchingToken == nil {
		return nil, errors.New("Refresh token not found or invalid")
	}

	// Check the expiration of the matching refresh token
	validityDuration := 7 * 24 * time.Hour
	if time.Now().Sub(matchingToken.CreatedAt) > validityDuration {
		// If the token has expired, return an error or handle it as needed
		return nil, errors.New("Refresh token has expired")
	}

	// If everything is valid, return the matching refresh token
	return matchingToken, nil
}
