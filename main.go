package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// By default gin serves on :8080 unless you specify a PORT in the environment
	// Make it run on port 8000
	r.Run("localhost:8000")
}
