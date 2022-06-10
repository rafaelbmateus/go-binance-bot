package binance

import (
	"fmt"
	"strings"

	"github.com/adshao/go-binance/v2"
)

// CreateOrder create new order in binance.
func (me *Binance) CreateOrder(symbol, buyWith, side string, quantity, price float64) (*binance.CreateOrderResponse, error) {
	order, err := me.Client.NewCreateOrderService().Symbol(fmt.Sprintf("%s%s", symbol, buyWith)).
		Side(getSide(side)).Type(binance.OrderTypeLimit).TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(toString(quantity)).Price(toString(price)).Do(*me.Context)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// getSide get which side type property by string.
func getSide(side string) binance.SideType {
	if strings.ToUpper(side) == "SELL" {
		return binance.SideTypeSell
	}

	return binance.SideTypeBuy
}

// toString converts a float to string with precision 5.
func toString(value float64) string {
	return fmt.Sprintf("%.5f", value)
}
