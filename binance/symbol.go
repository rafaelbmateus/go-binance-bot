package binance

import (
	"fmt"
	"math"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

// SymbolBalance get available amount in user wallet.
func (me *Binance) SymbolBalance(symbol string) (float64, error) {
	acc, err := me.Binance.NewGetAccountService().Do(*me.Context)
	if err != nil {
		return 0, err
	}

	for _, v := range acc.Balances {
		if v.Asset == symbol {
			free, err := strconv.ParseFloat(v.Free, 64)
			if err != nil {
				return 0, err
			}

			return math.Floor(free), nil
		}
	}

	return 0, nil
}

// SymbolPrice get the last symbol price.
func (me *Binance) SymbolPrice(symbol, buyWith, interval string) (float64, error) {
	klines, err := me.Binance.NewKlinesService().Symbol(fmt.Sprintf("%s%s", symbol, buyWith)).
		Interval(interval).Limit(1).Do(*me.Context)
	if err != nil {
		return 0, err
	}

	symbolPrice, err := strconv.ParseFloat(klines[0].Close, 64)
	if err != nil {
		return 0, err
	}

	return symbolPrice, nil
}

// SymbolExchangeInfo get the last symbol price.
func (me *Binance) SymbolExchangeInfo(symbol string) (*binance.ExchangeInfo, error) {
	res, err := me.Binance.NewExchangeInfoService().Symbol(symbol).Do(*me.Context)
	if err != nil {
		return nil, err
	}

	return res, nil
}
