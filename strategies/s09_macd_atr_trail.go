package strategies

import at "agenttrading"

// S09MACDATRTrail: S06's entry (MACD histogram > 0, close > EMA(50)), but exits only on a
// 3 ATR(14) chandelier stop from the highest close since entry or a close below EMA(50), so
// short histogram dips don't shake the trade out.
func S09MACDATRTrail(b at.Bars) []float64 {
	_, _, hist := at.MACD(b.Close, 12, 26, 9)
	trend, atr := at.EMA(b.Close, 50), at.ATR(b, 14)
	out := make([]float64, b.Len())
	held, peak := false, 0.0
	for i := range out {
		switch {
		case held:
			peak = max(peak, b.Close[i])
			held = b.Close[i] > peak-3*atr[i] && b.Close[i] > trend[i]
		case hist[i] > 0 && b.Close[i] > trend[i]:
			held, peak = true, b.Close[i]
		}
		if held {
			out[i] = 1
		}
	}
	return out
}
