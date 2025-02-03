package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rickj1ang/fly_crypto/api"
	"github.com/rickj1ang/fly_crypto/internal/app"
	"github.com/rickj1ang/fly_crypto/internal/checker"
	"github.com/rickj1ang/fly_crypto/internal/mail"
)

func initApp() (*app.App, error) {
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
	return app.NewApp(dbURI, redisURI)
}

// setupHealthCheck configures health check endpoints
func setupHealthCheck(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

// setupAuthRoutes configures authentication related routes
func setupAuthRoutes(r *gin.Engine, app *app.App) {
	r.POST("/login", api.Login(app))
	r.POST("/verify", api.Verify(app))
}

// setupProtectedRoutes configures routes that require authentication
func setupProtectedRoutes(r *gin.Engine, app *app.App) {
	protected := r.Group("/")
	protected.Use(api.AuthMiddleware(app))
	{
		protected.POST("/notifications", api.CreateNotification(app))
		protected.GET("/notifications", api.GetAllNotifications(app))
		protected.DELETE("/notifications/:id", api.DeleteNotification(app))
	}
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
	setupHealthCheck(r)
	setupAuthRoutes(r, app)
	setupProtectedRoutes(r, app)

	mailBox := make(chan mail.Message, 10)

	go mail.Sender(mailBox)
	checker.StartCheck(app, mailBox)

	// Start server
	if err := r.Run("localhost:8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
