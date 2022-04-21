package binance_test

import (
	"testing"

	"github.com/rafaelbmateus/binance-bot/binance"
	"github.com/stretchr/testify/assert"
)

func TestRSI(t *testing.T) {
	tests := []struct {
		name     string
		symbol   string
		interval string
		err      error
	}{
		{
			name:     "basic test",
			symbol:   "BNBUSDT",
			interval: "1d",
			err:      nil,
		},
	}
	b := binance.NewBinance(&log, &ctx, binanceClient)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsi, err := b.RSI(tt.symbol, tt.interval)
			assert.NotNil(t, rsi)
			assert.Equal(t, tt.err, err)
		})
	}
}
