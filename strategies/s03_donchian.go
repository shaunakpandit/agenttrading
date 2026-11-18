package strategies

import at "agenttrading"

// S03Donchian: long on a close above the prior 20-day high, short on a close below the prior
// 20-day low; each side exits on the opposite 10-day extreme.
func S03Donchian(b at.Bars) []float64 {
	high20, low20 := at.Shift(at.Highest(b.High, 20), 1), at.Shift(at.Lowest(b.Low, 20), 1)
	high10, low10 := at.Shift(at.Highest(b.High, 10), 1), at.Shift(at.Lowest(b.Low, 10), 1)
	out := make([]float64, b.Len())
	held := 0.0
	for i := range out {
		switch {
		case held > 0 && b.Close[i] < low10[i], held < 0 && b.Close[i] > high10[i]:
			held = 0
		}
		switch {
		case held == 0 && b.Close[i] > high20[i]:
			held = 1
		case held == 0 && b.Close[i] < low20[i]:
			held = -1
		}
		out[i] = held
	}
	return out
}
