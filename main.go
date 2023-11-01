package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		names := []string{"Yogendra Rana", "Ghanendra Rana"}
		c.HTML(200, "index.html", gin.H{
			"message": "Home page!",
			"names":   names,
		})
	})

	// By default gin serves on :8080 unless you specify a custom PORT by passing into the Run() method
	err := r.Run("localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
}
