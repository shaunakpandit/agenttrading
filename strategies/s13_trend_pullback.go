package strategies

import at "agenttrading"

// S13TrendPullback buys a dip inside an uptrend. With EMA(20) > EMA(50) and RSI(14) under 45 in the
// last 3 days, it enters when the close clears the prior day's high (the turn). It exits on a close
// below EMA(50) or a 3 ATR(14) chandelier stop.
func S13TrendPullback(b at.Bars) []float64 {
	fast, slow, atr := at.EMA(b.Close, 20), at.EMA(b.Close, 50), at.ATR(b, 14)
	rsi := at.RSI(b.Close, 14)
	recentDip := at.Lowest(rsi, 3)
	priorHigh := at.Shift(b.High, 1)
	out := make([]float64, b.Len())
	held, peak := false, 0.0
	for i := range out {
		switch {
		case held:
			peak = max(peak, b.Close[i])
			held = b.Close[i] > peak-3*atr[i] && b.Close[i] > slow[i]
		case fast[i] > slow[i] && recentDip[i] < 45 && b.Close[i] > priorHigh[i]:
			held, peak = true, b.Close[i]
		}
		if held {
			out[i] = 1
		}
	}
	return out
}
