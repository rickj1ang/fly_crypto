package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rickj1ang/fly_crypto/internal/app"
	baapi "github.com/rickj1ang/fly_crypto/internal/ba_api"
	"github.com/rickj1ang/fly_crypto/internal/data"
)

type createNotificationRequest struct {
	CoinSymbol  string  `json:"coin_symbol" binding:"required"`
	TargetPrice float64 `json:"target_price" binding:"required"`
}

// CreateNotification handles the creation of new price notifications
func CreateNotification(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createNotificationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert coin symbol to uppercase
		coinSymbol := strings.ToUpper(req.CoinSymbol)

		// Validate coin symbol
		allowedCoins := map[string]bool{"BTC": true, "ETH": true, "SOL": true}
		if !allowedCoins[coinSymbol] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coin symbol. Allowed values are: BTC, ETH, SOL"})
			return
		}

		// Get current price
		currentPrice, err := baapi.GetPrice(coinSymbol)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current price"})
			return
		}

		// Create notification using the data package struct
		notification := &data.Notification{
			CoinSymbol:  coinSymbol,
			TargetPrice: req.TargetPrice,
			IsAbove:     currentPrice < req.TargetPrice,
		}

		id, err := GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID from context"})
			return
		}

		notification.UserID = id
		// Save notification to the database
		err = a.Data.CreateNotification(notification)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "id": id})
			return
		}
		email, err := a.Data.GetUserEmail(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user email"})
			return
		}

		key := GetKey(notification)
		if err = a.Data.StoreInSortedSet(key, req.TargetPrice, email); err != nil {
			//delete notification from db
			if err = a.Data.DeleteNotification(notification.ID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store notification in cache"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Notification created successfully",
			"data":    notification,
		})
	}
}
