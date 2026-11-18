package strategies

import at "agenttrading"

// S05BreakoutATRTrail: long on a close above the prior 20-day high; exit when the close falls
// 3 ATR(14) below the highest close since entry (chandelier stop on closes).
func S05BreakoutATRTrail(b at.Bars) []float64 {
	high20, atr := at.Shift(at.Highest(b.High, 20), 1), at.ATR(b, 14)
	out := make([]float64, b.Len())
	held, peak := false, 0.0
	for i := range out {
		switch {
		case held:
			peak = max(peak, b.Close[i])
			held = b.Close[i] > peak-3*atr[i]
		case b.Close[i] > high20[i]:
			held, peak = true, b.Close[i]
		}
		if held {
			out[i] = 1
		}
	}
	return out
}
