package binance

import (
	"math"
	"strconv"
)

func (me *Binance) RSI(symbol, interval string) (float64, error) {
	klines, err := me.Binance.NewKlinesService().Symbol(symbol).
		Interval(interval).Limit(1).Do(*me.Context)
	if err != nil {
		return 0, err
	}

	totalGain := 0.0
	totalLoss := 0.0

	for i := 1; i < len(klines); i++ {
		previous := klines[i].Close
		current := klines[i-1].Close

		// convert string to float64
		previousClose, _ := strconv.ParseFloat(previous, 64)
		currentClose, _ := strconv.ParseFloat(current, 64)

		difference := currentClose - previousClose
		if difference >= 0 {
			totalGain += difference
		} else {
			totalLoss -= difference
		}
	}

	rs := totalGain / math.Abs(totalLoss)
	rsi := 100 - 100/(1+rs)

	me.Log.Debug().Msgf("[%s] RSI: %f", symbol, rsi)
	return rsi, nil
}
