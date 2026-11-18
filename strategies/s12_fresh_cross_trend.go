package strategies

import at "agenttrading"

// S12FreshCrossTrend: enter only on the bar the MACD(12, 26, 9) histogram crosses above zero
// with the close above EMA(50). Exit on a 3 ATR(14) chandelier stop or a close below EMA(50).
// No re-entry until the next fresh cross (the S09 whipsaw fix).
func S12FreshCrossTrend(b at.Bars) []float64 {
	_, _, hist := at.MACD(b.Close, 12, 26, 9)
	zero := make([]float64, b.Len())
	cross := at.CrossOver(hist, zero)
	trend, atr := at.EMA(b.Close, 50), at.ATR(b, 14)
	out := make([]float64, b.Len())
	held, peak := false, 0.0
	for i := range out {
		switch {
		case held:
			peak = max(peak, b.Close[i])
			held = b.Close[i] > peak-3*atr[i] && b.Close[i] > trend[i]
		case cross[i] && b.Close[i] > trend[i]:
			held, peak = true, b.Close[i]
		}
		if held {
			out[i] = 1
		}
	}
	return out
}
