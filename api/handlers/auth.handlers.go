package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User

func RegisterHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	for _, u := range users {
		if u.Username == user.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}
	}

	// Add user to list of users
	users = append(users, user)

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func LoginHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	for _, u := range users {
		if u.Username == user.Username && u.Password == user.Password {
			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}
