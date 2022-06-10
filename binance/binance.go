package binance

import (
	"context"

	"github.com/adshao/go-binance/v2"
	"github.com/rs/zerolog"
)

// Binance represents this package.
type Binance struct {
	Log     *zerolog.Logger
	Context *context.Context
	Client  *binance.Client
}

// NewBinance create a new binance instance.
func NewBinance(log *zerolog.Logger, ctx *context.Context, binance *binance.Client) *Binance {
	return &Binance{
		Log:     log,
		Context: ctx,
		Client:  binance,
	}
}
