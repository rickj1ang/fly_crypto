package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rickj1ang/fly_crypto/internal/app"
)

// AuthMiddleware verifies the bearer token in the Authorization header
func AuthMiddleware(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			return
		}

		tokenStr := parts[1]

		// Verify token exists in Redis and get associated email
		userID, err := a.GetUserIDByAuthToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Set user_id in context for later use
		c.Set("user_id", userID)
		c.Next()
	}
}
