package strategies

import at "agenttrading"

// S07SqueezeBreakout: when Bollinger (20, 2) width hit its 120-day low within the last 10 days and
// the close breaks above the upper band, go long; exit on a close below EMA(20).
func S07SqueezeBreakout(b at.Bars) []float64 {
	mid, upper, lower := at.Bollinger(b.Close, 20, 2)
	width := make([]float64, b.Len())
	for i := range width {
		width[i] = (upper[i] - lower[i]) / mid[i]
	}
	recentSqueeze := at.Lowest(width, 120)
	exitLine := at.EMA(b.Close, 20)
	out := make([]float64, b.Len())
	held, lastSqueeze := false, -1000
	for i := range out {
		if width[i] <= recentSqueeze[i] {
			lastSqueeze = i
		}
		switch {
		case held:
			held = b.Close[i] >= exitLine[i]
		case i-lastSqueeze <= 10 && b.Close[i] > upper[i]:
			held = true
		}
		if held {
			out[i] = 1
		}
	}
	return out
}
