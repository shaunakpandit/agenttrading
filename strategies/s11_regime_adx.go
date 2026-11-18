package strategies

import at "agenttrading"

// S11RegimeADX switches rules by regime. When ADX(14) >= 20 (trending), it is long while
// EMA(20) > EMA(50) and the close is above EMA(50). Below 20 (ranging), it buys closes under the
// lower Bollinger band (20, 2) and sells at the middle band.
func S11RegimeADX(b at.Bars) []float64 {
	adx := at.ADX(b, 14)
	fast, slow := at.EMA(b.Close, 20), at.EMA(b.Close, 50)
	mid, _, lower := at.Bollinger(b.Close, 20, 2)
	out := make([]float64, b.Len())
	held := false
	for i := range out {
		if adx[i] >= 20 {
			held = fast[i] > slow[i] && b.Close[i] > slow[i]
		} else if held {
			held = b.Close[i] < mid[i]
		} else {
			held = b.Close[i] < lower[i]
		}
		if held {
			out[i] = 1
		}
	}
	return out
}
