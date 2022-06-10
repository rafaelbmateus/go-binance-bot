package strategy

import (
	"strconv"

	"github.com/adshao/go-binance/v2"
)

func CalculateRSI(klines []*binance.Kline) (float64, error) {
	totalGain := 0.0
	totalLoss := 0.0
	for i := 1; i < len(klines); i++ {
		current, _ := strconv.ParseFloat(klines[i].Close, 64)
		previous, _ := strconv.ParseFloat(klines[i-1].Close, 64)
		diff := current - previous
		if diff >= 0 {
			totalGain += diff
		} else {
			totalLoss -= diff
		}
	}

	avgGain := totalGain / float64(len(klines))
	avgLoss := totalLoss / float64(len(klines))
	rs := avgGain / avgLoss
	rsi := 100 - 100/(rs+1)

	return rsi, nil
}
