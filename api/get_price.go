package api

import (
	"encoding/json"
	"net/http"
)

const baseURL = "https://data-api.binance.vision/api/v3/ticker/price?symbol="

func GetPrice(coin string) (string, error) {
	res, err := http.Get(getURL(coin))
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}
	res.Body.Close()

	price := result["price"]
	return price.(string), nil
}

func getURL(coin string) string {
	return baseURL + coin + "USDT"
}
