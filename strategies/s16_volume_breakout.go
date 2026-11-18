package strategies

import at "agenttrading"

// S16VolumeBreakout: S05 (20-day closing breakout, 3 ATR chandelier exit), but the breakout day's
// volume must be at least 1.5x its 20-day average, so thin breakouts are skipped.
func S16VolumeBreakout(b at.Bars) []float64 {
	high20, atr := at.Shift(at.Highest(b.High, 20), 1), at.ATR(b, 14)
	averageVolume := at.Shift(at.SMA(b.Volume, 20), 1)
	out := make([]float64, b.Len())
	held, peak := false, 0.0
	for i := range out {
		switch {
		case held:
			peak = max(peak, b.Close[i])
			held = b.Close[i] > peak-3*atr[i]
		case b.Close[i] > high20[i] && b.Volume[i] > 1.5*averageVolume[i]:
			held, peak = true, b.Close[i]
		}
		if held {
			out[i] = 1
		}
	}
	return out
}
