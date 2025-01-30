package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rickj1ang/fly_crypto/internal/app"
)

func GetAllNotifications(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Get all notifications for the user
		notifications, err := a.Data.GetUserAllNotifications(userID.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"notifications": notifications})
	}
}
