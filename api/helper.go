package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rickj1ang/fly_crypto/internal/data"
)

func GetUserIDFromContext(c *gin.Context) (int64, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("user ID not found in context")
	}

	// Type assertion
	userIDInt64, ok := userID.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID type in context")
	}

	return userIDInt64, nil
}

// GetKey generates a Redis key for storing notification in sorted set
func GetKey(notification *data.Notification) string {
	var direction string
	if notification.IsAbove {
		direction = "above"
	} else {
		direction = "below"
	}
	return fmt.Sprintf("%s:%s", notification.CoinSymbol, direction)
}
