package strategies

import at "agenttrading"

// S01TrendSMA: long while close > SMA(50) and SMA(20) > SMA(50); flat otherwise.
func S01TrendSMA(b at.Bars) []float64 {
	fast, slow := at.SMA(b.Close, 20), at.SMA(b.Close, 50)
	out := make([]float64, b.Len())
	for i := range out {
		if b.Close[i] > slow[i] && fast[i] > slow[i] {
			out[i] = 1
		}
	}
	return out
}
