package handlers

import (
	"errors"
	"go-gin-auth/src/models"
	"go-gin-auth/src/services"
	"go-gin-auth/src/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	bcryptCost = 10
)

type RegisterInput struct {
	Name            string `json:"name" binding:"required,min=3,max=50"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// sign up controller
func Register(ctx *gin.Context) {
	var input RegisterInput

	// get db connection
	db := ctx.MustGet("db").(*gorm.DB)

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// validation checks
	if input.Password != input.ConfirmPassword {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": "Passwords do not match"})
		return
	}

	// check if user already exists
	result := db.Where("email = ?", input.Email).First(&models.User{})
	if result.Error == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": "User already exists."})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcryptCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash password"})
		return
	}

	// create user
	user := models.User{Name: input.Name, Email: input.Email, Password: string(hashedPassword)}
	result = db.Create(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create user"})
		return
	}

	// generate access tokens
	signedAccessToken, err := services.GenerateAccessToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}

	// generate refresh token
	refreshToken, hashedToken, err := services.GenerateRefreshTokenAndHash()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}

	// Save refresh token in the database
	dbRefreshToken := models.RefreshToken{UserID: user.ID, TokenHash: hashedToken}
	result = db.Create(&dbRefreshToken)
	if result.Error != nil || result.RowsAffected == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save refresh token"})
		return
	}

	// set http only cookies (name, value, maxAge, path, domain, secure, httpOnly)
	ctx.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "User created successfully",
		Data: gin.H{
			"access_token": signedAccessToken,
			"user":         user,
		},
	})
}

// sign in controller
func Login(ctx *gin.Context) {
	var input LoginInput

	// get db connection
	db := ctx.MustGet("db").(*gorm.DB)

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if user already exists
	var user models.User
	result := db.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	// compare hash password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// generate access tokens
	signedAccessToken, err := services.GenerateAccessToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// generate refresh token
	refreshToken, hashedToken, err := services.GenerateRefreshTokenAndHash()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// update refresh token in the database
	dbRefreshToken := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashedToken,
	}

	updateFields := map[string]interface{}{
		"token_hash": dbRefreshToken.TokenHash,
	}

	result = db.Model(&models.RefreshToken{}).Where("user_id = ?", user.ID).Updates(updateFields)

	// set http only cookies (name, value, maxAge, path, domain, secure, httpOnly)
	ctx.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "User logged in successfully",
		Data: gin.H{
			"access_token": signedAccessToken,
			"user":         user,
		},
	})
}

// delete user
func Logout(c *gin.Context) {

}
