package strategies

import (
	"math"

	at "agenttrading"
)

// S17EfficientBreakout: S05, but only enter when Kaufman's efficiency ratio over 20 days
// (|net move| / sum of |daily moves|) is above 0.3, i.e. the stock has been trending cleanly
// rather than chopping.
func S17EfficientBreakout(b at.Bars) []float64 {
	high20, atr := at.Shift(at.Highest(b.High, 20), 1), at.ATR(b, 14)
	out := make([]float64, b.Len())
	held, peak := false, 0.0
	for i := range out {
		switch {
		case held:
			peak = max(peak, b.Close[i])
			held = b.Close[i] > peak-3*atr[i]
		case b.Close[i] > high20[i] && efficiencyRatio(b.Close, i, 20) > 0.3:
			held, peak = true, b.Close[i]
		}
		if held {
			out[i] = 1
		}
	}
	return out
}

func efficiencyRatio(close []float64, i, n int) float64 {
	if i < n {
		return math.NaN()
	}
	path := 0.0
	for j := i - n + 1; j <= i; j++ {
		path += math.Abs(close[j] - close[j-1])
	}
	return math.Abs(close[i]-close[i-n]) / path
}
