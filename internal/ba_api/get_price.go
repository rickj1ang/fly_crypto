package baapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const baseURL = "https://data-api.binance.vision/api/v3/ticker/price?symbol="

func GetPrice(coin string, prices *sync.Map) float64 {
	price, _ := prices.Load(coin)
	return price.(float64)
}

func fetchPrice(coin string) (float64, error) {
	res, err := http.Get(getURL(coin))
	if err != nil {
		return -1, err
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return -1, err
	}
	res.Body.Close()

	price := result["price"]

	currentPrice, err := strconv.ParseFloat(price.(string), 64)
	return currentPrice, nil
}

func getURL(coin string) string {
	return baseURL + coin + "USDT"
}

func PriceUpdater(coins []string, prices *sync.Map) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for _, coin := range coins {
				price, err := fetchPrice(coin)
				if err != nil {
					continue
				}
				prices.Store(coin, price)
			}
		}
	}
}

func InitPrice(coins []string, prices *sync.Map) {
	for _, coin := range coins {
		price, err := fetchPrice(coin)
		if err != nil {
			panic(err)
		}
		prices.Store(coin, price)
	}

}
