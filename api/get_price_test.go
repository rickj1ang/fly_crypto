package api

import (
	"testing"
)

func TestGetPrice(t *testing.T) {
	coins := []string{"BTC", "ETH", "BNB", "SOL", "AR"}

	for _, coin := range coins {
		price, err := GetPrice(coin)
		if err != nil {
			t.Errorf("fail to get %s price: %s", coin, err.Error())
		}
		if price == "" {
			t.Errorf("fail to get %s price", coin)
		}
		t.Log(price)
	}

}
