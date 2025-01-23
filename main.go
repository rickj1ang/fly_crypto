package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rickj1ang/fly_crypto/api"
)

func initApp() (*api.App, error) {
	// Get database connection strings from environment variables
	dbURI := os.Getenv("DATABASE_URL")
	if dbURI == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	redisURI := os.Getenv("REDIS_URL")
	if redisURI == "" {
		log.Fatal("REDIS_URL environment variable is not set")
	}

	// Initialize app with database connections
	return api.NewApp(dbURI, redisURI)
}

func main() {
	// Initialize app with database connections
	app, err := initApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer app.Close()

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Login routes
	r.POST("/login", app.Login())
	r.POST("/verify", app.Verify())

	// Protected routes group
	protected := r.Group("/")
	protected.Use(app.AuthMiddleware())
	{
		// Add protected routes here
	}

	// Start server
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
