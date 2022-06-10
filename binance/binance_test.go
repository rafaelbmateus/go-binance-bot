package binance_test

import (
	"context"
	"os"
	"testing"

	sdk "github.com/adshao/go-binance/v2"
	"github.com/rafaelbmateus/binance-bot/binance"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	ctx           = context.Background()
	log           = zerolog.New(os.Stdout)
	apiKey        = os.Getenv("BINANCE_API_KEY")
	apiSecret     = os.Getenv("BINANCE_API_SECRET")
	binanceClient = sdk.NewClient(apiKey, apiSecret)
)

func TestNewBinance(t *testing.T) {
	b := binance.NewBinance(&log, &ctx, binanceClient)
	assert.Equal(t, "https://api.binance.com", b.Client.BaseURL)
	assert.NotNil(t, b.Context)
	assert.NotNil(t, b.Log)
}
