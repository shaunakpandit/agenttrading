package strategies

import at "agenttrading"

// S15TimeSeriesMomentum: long while the close is above its level 63 days ago (3-month momentum)
// and above EMA(20) (short-term confirmation).
func S15TimeSeriesMomentum(b at.Bars) []float64 {
	ago, fast := at.Shift(b.Close, 63), at.EMA(b.Close, 20)
	out := make([]float64, b.Len())
	for i := range out {
		if b.Close[i] > ago[i] && b.Close[i] > fast[i] {
			out[i] = 1
		}
	}
	return out
}
