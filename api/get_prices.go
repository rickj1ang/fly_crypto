package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rickj1ang/fly_crypto/internal/app"
	baapi "github.com/rickj1ang/fly_crypto/internal/ba_api"
)

// GetPrices returns current prices for all supported cryptocurrencies
func GetPrices(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		prices := make(map[string]float64)

		// Get price for each supported coin
		for _, coin := range a.SupportCoins {
			price := baapi.GetPrice(coin, a.CoinsPrices)
			prices[coin] = price
		}

		// Return prices as JSON
		c.JSON(http.StatusOK, gin.H{
			"prices": prices,
		})

	}
}
