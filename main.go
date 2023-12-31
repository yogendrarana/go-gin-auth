package main

import (
	"go-gin-auth/src/initializers"
	"go-gin-auth/src/middlewares"
	"go-gin-auth/src/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
	initializers.SyncDatabase()
}

func main() {
	router := gin.New()
	router.LoadHTMLGlob("./src/views/*")

	// Middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gin.ErrorLogger())
	router.Use(middlewares.ErrorMiddleware)

	// Serve static files
	router.GET("/", func(c *gin.Context) { c.File("./src/views/index.html") })

	// Initialize routes
	apiV1 := router.Group("/api/v1")
	routes.AuthRoutes(apiV1)
	routes.ProtectedRoutes(apiV1)

	// By default gin serves on :8080 unless you specify a custom PORT by passing into the Run() method
	err := router.Run("localhost:" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
