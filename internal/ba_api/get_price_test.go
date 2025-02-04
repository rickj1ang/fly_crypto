package baapi

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPriceUpdater(t *testing.T) {
	coins := []string{
		"BTC",
		"ETH",
		"SOL",
	}

	prices := &sync.Map{}
	go PriceUpdater(coins, prices)
	for {
		price, ok := prices.Load("BTC")
		if ok {
			fmt.Println(price)
		}
		time.Sleep(time.Second)
	}
}
