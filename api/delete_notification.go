package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rickj1ang/fly_crypto/internal/app"
)

// DeleteNotification handles the deletion of a notification
func DeleteNotification(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get notification ID from URL parameter
		notificationIDStr := c.Param("notification_id")
		notificationID, err := strconv.ParseInt(notificationIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
			return
		}
		notification, err := a.Data.GetANotification(notificationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notification"})
			return
		}

		// Get user ID from context
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Check if the notification belongs to the user
		notificationIDs, err := a.Data.GetUserNotificationIDs(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify notification ownership"})
			return
		}

		belongsToUser := false
		for _, id := range notificationIDs {
			if id == notificationID {
				belongsToUser = true
				break
			}
		}

		if !belongsToUser {
			c.JSON(http.StatusForbidden, gin.H{"error": "Notification does not belong to user"})
			return
		}

		// Delete the notification
		if err := a.Data.DeleteNotification(notificationID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
			return
		}

		//remove the notification from redis
		key := GetKey(notification)
		email, err := a.Data.GetUserEmail(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user email"})
			return
		}
		if err = a.Data.RemoveFromSortedSet(key, email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification from redis"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
	}
}
