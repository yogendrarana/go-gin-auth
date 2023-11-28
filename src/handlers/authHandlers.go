package handlers

import (
	"errors"
	custom_errors "go-gin-auth/src/errors"
	"go-gin-auth/src/models"
	"go-gin-auth/src/services"
	"go-gin-auth/src/types"
	"net/http"
	"time"

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
func Register(c *gin.Context) {
	var input RegisterInput

	// get db connection
	db := c.MustGet("db").(*gorm.DB)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// validation checks
	if input.Password != input.ConfirmPassword {
		err := custom_errors.NewAppError("Password do not match.", 400)
		c.Error(err)
		return
	}

	// check if user already exists
	result := db.Where("email = ?", input.Email).First(&models.User{})
	if result.Error == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"success": false, "message": "User already exists."})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcryptCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash password"})
		return
	}

	// create user
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	result = db.Create(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create user"})
		return
	}

	// generate access tokens
	accessSignedToken, err := services.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}

	// generate refresh token
	refreshToken, hashedToken, err := services.GenerateRefreshTokenAndHash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}

	// Save refresh token in the database
	dbRefreshToken := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	result = db.Create(&dbRefreshToken)
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save refresh token"})
		return
	}

	// set http only cookies (name, value, maxAge, path, domain, secure, httpOnly)
	c.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "localhost", false, true)

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "User created successfully",
		Data: gin.H{
			"access_token": accessSignedToken,
			"user":         user,
		},
	})
}

// sign in controller
func Login(c *gin.Context) {
	var input LoginInput

	// get db connection
	db := c.MustGet("db").(*gorm.DB)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if user already exists
	var user models.User
	result := db.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	// compare hash password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// generate access tokens
	accessSignedToken, err := services.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// generate refresh token
	refreshToken, hashedToken, err := services.GenerateRefreshTokenAndHash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// update refresh token in the database
	dbRefreshToken := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	updateFields := map[string]interface{}{
		"token_hash": dbRefreshToken.TokenHash,
		"expires_at": dbRefreshToken.ExpiresAt,
	}

	result = db.Model(&models.RefreshToken{}).Where("user_id = ?", user.ID).Updates(updateFields)

	// set http only cookies (name, value, maxAge, path, domain, secure, httpOnly)
	c.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "localhost", false, true)

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Message: "User logged in successfully",
		Data: gin.H{
			"access_token": accessSignedToken,
			"user":         user,
		},
	})
}

// delete user
func DeleteUser(c *gin.Context) {

}
