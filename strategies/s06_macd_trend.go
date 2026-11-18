package strategies

import at "agenttrading"

// S06MACDTrend: long while the MACD(12, 26, 9) histogram is positive and the close is above EMA(50).
func S06MACDTrend(b at.Bars) []float64 {
	_, _, hist := at.MACD(b.Close, 12, 26, 9)
	trend := at.EMA(b.Close, 50)
	out := make([]float64, b.Len())
	for i := range out {
		if hist[i] > 0 && b.Close[i] > trend[i] {
			out[i] = 1
		}
	}
	return out
}
