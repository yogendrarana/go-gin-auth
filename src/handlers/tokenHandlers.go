package handlers

import (
	"go-gin-auth/src/models"
	"go-gin-auth/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RefreshTokenInput struct {
	UserId uint `json:"user_id" binding:"required"`
}

func HandleRefreshToken(ctx *gin.Context) {
	var input RefreshTokenInput

	// get db connection
	db := ctx.MustGet("db").(*gorm.DB)

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// retrieve refresh token from cookie
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"success": false, "message": "Refresh token not found in cookie."})
		return
	}

	// check if refresh token is present in database
	refreshTokenArrayPtr, err := models.GetRefreshTokenByUserID(db, input.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"success": false, "message": "Refresh token not found in database."})
		return
	}

	refreshTokenArray := *refreshTokenArrayPtr

	// validate refresh token in cookie
	validatedRefreshToken, err := services.ValidateRefreshToken(refreshToken, refreshTokenArray)

	if validatedRefreshToken == nil || err != nil {
		// cookie syntax: SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool)
		ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

		ctx.AbortWithStatusJSON(400, gin.H{"success": false, "message": "Refresh token is invalid. Please login again."})
		return
	}

	// generate access tokens
	signedAccessToken, err := services.GenerateAccessToken(validatedRefreshToken.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}

	// send the access token in the response
	ctx.AbortWithStatusJSON(200, gin.H{
		"success": true,
		"message": "Access token generated successfully.",
		"data": gin.H{
			"access_token": signedAccessToken,
		},
	})

}
