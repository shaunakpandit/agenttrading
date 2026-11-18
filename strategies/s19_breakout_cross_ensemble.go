package strategies

import at "agenttrading"

// S19BreakoutCrossEnsemble: half S05 (price breakout), half S12 (MACD fresh cross). Each was
// the most consistent of its kind, and they enter on different evidence.
func S19BreakoutCrossEnsemble(b at.Bars) []float64 {
	return average(S05BreakoutATRTrail(b), S12FreshCrossTrend(b))
}

// average is the bar-by-bar mean of several position series.
func average(series ...[]float64) []float64 {
	out := make([]float64, len(series[0]))
	for _, s := range series {
		for i := range out {
			out[i] += s[i] / float64(len(series))
		}
	}
	return out
}
