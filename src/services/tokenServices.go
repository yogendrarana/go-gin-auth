package services

import (
	"os"
	"time"

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
