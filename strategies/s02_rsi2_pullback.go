package strategies

import at "agenttrading"

// S02RSI2Pullback: in an uptrend (close > SMA(100)) buy when RSI(2) < 10; sell when the close
// rises above SMA(5) or after 10 days.
func S02RSI2Pullback(b at.Bars) []float64 {
	rsi, trend, exitLine := at.RSI(b.Close, 2), at.SMA(b.Close, 100), at.SMA(b.Close, 5)
	out := make([]float64, b.Len())
	held := 0
	for i := range out {
		switch {
		case held > 0 && (b.Close[i] > exitLine[i] || held >= 10):
			held = 0
		case held > 0:
			held++
		case rsi[i] < 10 && b.Close[i] > trend[i]:
			held = 1
		}
		if held > 0 {
			out[i] = 1
		}
	}
	return out
}
