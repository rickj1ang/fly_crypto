package baapi

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const baseURL = "https://data-api.binance.vision/api/v3/ticker/price?symbol="

func GetPrice(coin string) (float64, error) {
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
