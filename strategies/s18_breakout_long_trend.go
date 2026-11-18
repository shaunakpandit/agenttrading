package strategies

import at "agenttrading"

// S18BreakoutLongTrend: S05, but only enter while the close is above SMA(150), so breakouts
// inside a long-term downtrend are skipped. The exit is unchanged.
func S18BreakoutLongTrend(b at.Bars) []float64 {
	high20, atr, longTrend := at.Shift(at.Highest(b.High, 20), 1), at.ATR(b, 14), at.SMA(b.Close, 150)
	out := make([]float64, b.Len())
	held, peak := false, 0.0
	for i := range out {
		switch {
		case held:
			peak = max(peak, b.Close[i])
			held = b.Close[i] > peak-3*atr[i]
		case b.Close[i] > high20[i] && b.Close[i] > longTrend[i]:
			held, peak = true, b.Close[i]
		}
		if held {
			out[i] = 1
		}
	}
	return out
}
