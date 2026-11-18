package strategies

import at "agenttrading"

// S20TrendReversionEnsemble: half S05 (trend breakout), half S14 (IBS dip buying). Trend and
// reversion tend to win in different regimes, so their sum should be smoother than either.
func S20TrendReversionEnsemble(b at.Bars) []float64 {
	return average(S05BreakoutATRTrail(b), S14IBSReversion(b))
}
