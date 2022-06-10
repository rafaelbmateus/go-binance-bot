package bot_test

import (
	"context"
	"os"
	"testing"

	sdk "github.com/adshao/go-binance/v2"
	"github.com/rafaelbmateus/binance-bot/binance"
	"github.com/rafaelbmateus/binance-bot/bot"
	"github.com/rafaelbmateus/binance-bot/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	log       = zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx       = context.Background()
	apiKey    = os.Getenv("BINANCE_API_KEY")
	apiSecret = os.Getenv("BINANCE_API_SECRET")
	trade     = config.Trade{
		Symbol:      "BTC/USDT",
		Amount:      0,
		RSIBuy:      40,
		RSISell:     60,
		RSIInterval: "5m",
		RSILimit:    14,
	}
)

func TestMonitor(t *testing.T) {
	binance := binance.NewBinance(&log, &ctx, sdk.NewClient(apiKey, apiSecret))
	bot := bot.NewBot(&log, &ctx, binance, nil, nil)
	err := bot.Monitor(trade)
	assert.NoError(t, err)
}
