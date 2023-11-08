package utils

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateJWTTokensForUser(userID uint) (string, string, error) {
	// access token expires in 15 minutes and refresh token expires in 24 hours
	accessExpiry := time.Now().Add(time.Minute * 15).Unix()
	refreshExpiry := time.Now().Add(time.Hour * 24).Unix()

	// defining the claims for the tokens
	accessClaims := jwt.MapClaims{"exp": accessExpiry, "sub": userID}
	refreshClaims := jwt.MapClaims{"exp": refreshExpiry, "sub": userID}

	// create the access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSignedToken, err := accessToken.SignedString([]byte(os.Getenv("ACCESS_JWT_SECRET")))
	if err != nil {
		return "", "", err
	}

	// create the refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshSignedToken, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_JWT_SECRET")))
	if err != nil {
		return "", "", err
	}

	return accessSignedToken, refreshSignedToken, nil
}
