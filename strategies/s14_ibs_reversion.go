package strategies

import at "agenttrading"

// S14IBSReversion: internal bar strength (close - low) / (high - low) below 0.15 means the stock
// closed near its low. Buy that above SMA(100); sell on a close above the prior day's high or
// after 5 days.
func S14IBSReversion(b at.Bars) []float64 {
	trend, priorHigh := at.SMA(b.Close, 100), at.Shift(b.High, 1)
	out := make([]float64, b.Len())
	held := 0
	for i := range out {
		ibs := (b.Close[i] - b.Low[i]) / (b.High[i] - b.Low[i])
		switch {
		case held > 0 && (b.Close[i] > priorHigh[i] || held >= 5):
			held = 0
		case held > 0:
			held++
		case ibs < 0.15 && b.Close[i] > trend[i]:
			held = 1
		}
		if held > 0 {
			out[i] = 1
		}
	}
	return out
}
