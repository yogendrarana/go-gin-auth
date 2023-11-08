package main

import (
	"go-gin-auth/src/initializers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
}

func main() {
	router := gin.New()
	router.LoadHTMLGlob("./src/views/*")

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gin.ErrorLogger())

	router.GET("/", func(c *gin.Context) {
		c.File("./src/views/index.html")
	})

	// By default gin serves on :8080 unless you specify a custom PORT by passing into the Run() method
	err := router.Run("localhost:" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
