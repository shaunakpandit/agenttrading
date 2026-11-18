package strategies

import at "agenttrading"

// S04BollingerReversion: buy a close below the lower Bollinger band (20, 2); sell at the middle
// band, or after 15 days.
func S04BollingerReversion(b at.Bars) []float64 {
	mid, _, lower := at.Bollinger(b.Close, 20, 2)
	out := make([]float64, b.Len())
	held := 0
	for i := range out {
		switch {
		case held > 0 && (b.Close[i] > mid[i] || held >= 15):
			held = 0
		case held > 0:
			held++
		case b.Close[i] < lower[i]:
			held = 1
		}
		if held > 0 {
			out[i] = 1
		}
	}
	return out
}
