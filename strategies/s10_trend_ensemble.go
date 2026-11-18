package strategies

import at "agenttrading"

// S10TrendEnsemble: the average of three different trend rules (S01 moving averages, S05
// breakout with trailing stop, S06 MACD), so size scales with how many agree.
func S10TrendEnsemble(b at.Bars) []float64 {
	votes := [][]float64{S01TrendSMA(b), S05BreakoutATRTrail(b), S06MACDTrend(b)}
	out := make([]float64, b.Len())
	for _, v := range votes {
		for i := range out {
			out[i] += v[i] / float64(len(votes))
		}
	}
	return out
}
