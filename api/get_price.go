package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://data-api.binance.vision/api/v3/ticker/price?symbol="

func GetPrice(coin string) string {
	res, err := http.Get(getURL(coin))
	if err != nil {
		fmt.Println(err)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		panic(err)
	}
	res.Body.Close()

	price := result["price"]
	return price.(string)
}

func getURL(coin string) string {
	return baseURL + coin + "USDT"
}
