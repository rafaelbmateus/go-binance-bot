package strategy_test

import (
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/rafaelbmateus/binance-bot/strategy"
	"github.com/stretchr/testify/assert"
)

func TestCalculateRSI(t *testing.T) {
	tests := []struct {
		name   string
		klines []*binance.Kline
		rsi    float64
		err    error
	}{
		{
			name: "basic test",
			klines: []*binance.Kline{
				{Close: "10"},
				{Close: "9"},
				{Close: "12"},
				{Close: "5"},
				{Close: "6"},
				{Close: "6"},
				{Close: "6"},
				{Close: "8"},
				{Close: "9"},
				{Close: "10"},
				{Close: "12"},
				{Close: "13"},
				{Close: "14"},
				{Close: "15"},
			},
			rsi: 61.904761904761905,
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsi, err := strategy.CalculateRSI(tt.klines)
			assert.Equal(t, tt.rsi, rsi)
			assert.Equal(t, tt.err, err)
		})
	}
}
