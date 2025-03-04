package checker

import (
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/rickj1ang/fly_crypto/internal/app"
	baapi "github.com/rickj1ang/fly_crypto/internal/ba_api"
	"github.com/rickj1ang/fly_crypto/internal/mail"
)

func isAbove(targetPrice float64, coinSymbol string, prices *sync.Map) bool {
	currentPrice := baapi.GetPrice(coinSymbol, prices)

	return targetPrice <= currentPrice
}

func isBelow(targetPrice float64, coinSymbol string, prices *sync.Map) bool {
	currentPrice := baapi.GetPrice(coinSymbol, prices)

	return targetPrice >= currentPrice
}

func checkAbove(app *app.App, coinSymbol string) bool {
	key := coinSymbol + ":above"
	minFromAbove, _, err := app.Data.GetMinScoreFromSortedSet(key)
	if err != nil {
		if err != redis.Nil {
			panic(err)
		}
		return false
	}
	res := isAbove(minFromAbove, coinSymbol, app.CoinsPrices)

	return res
}

func checkBelow(app *app.App, coinSymbol string) bool {
	key := coinSymbol + ":below"
	maxFromAbove, _, err := app.Data.GetMaxScoreFromSortedSet(key)
	if err != nil {
		if err != redis.Nil {
			panic(err)
		}
		return false
	}
	res := isBelow(maxFromAbove, coinSymbol, app.CoinsPrices)

	return res
}

func sendAbove(app *app.App, coinSymbol string, box chan mail.Message) error {
	key := coinSymbol + ":above"
	price, email, err := app.Data.PopMinFromSortedSet(key)
	if err != nil {
		return err
	}
	message := mail.Message{
		SendTo:      email.(string),
		TargetPrice: price,
		CoinSymbol:  coinSymbol,
	}
	box <- message
	err = app.Data.DeleteNotificationFromMessage(message.SendTo, true, message.CoinSymbol)
	if err != nil {
		panic(err)
	}
	return nil
}

func sendBelow(app *app.App, coinSymbol string, box chan mail.Message) error {
	key := coinSymbol + ":below"
	price, email, err := app.Data.PopMaxFromSortedSet(key)
	if err != nil {
		return err
	}
	message := mail.Message{
		SendTo:      email.(string),
		TargetPrice: price,
		CoinSymbol:  coinSymbol,
	}
	box <- message
	err = app.Data.DeleteNotificationFromMessage(message.SendTo, false, message.CoinSymbol)
	if err != nil {
		panic(err)
	}
	return nil
}

func StartCheck(app *app.App, box chan mail.Message) {
	for _, coin := range app.SupportCoins {
		go checkOneCoin(app, coin, box)
	}
}

func checkOneCoin(app *app.App, coinSymbol string, box chan mail.Message) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if checkAbove(app, coinSymbol) {
				sendAbove(app, coinSymbol, box)
			}
			if checkBelow(app, coinSymbol) {
				sendBelow(app, coinSymbol, box)
			}
		}
	}
}
