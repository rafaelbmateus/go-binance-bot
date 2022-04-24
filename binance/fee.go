package binance

import (
	"strings"

	"github.com/rafaelbmateus/binance-bot/config"
)

// Fee check trading fee rate.
func Fee(symbol string) float64 {
	if strings.ToUpper(symbol) == "BNB" {
		return 0.075
	}

	return 0.1
}

func ProfitPerc(trade config.Trade) float64 {
	return (trade.SellPrice - trade.BuyPrice) / (trade.SellPrice) * 100
}

func Profit(trade config.Trade) float64 {
	return (Fee(trade.BuyWith()) / 100) * (trade.BuyPrice + trade.SellPrice)
}
