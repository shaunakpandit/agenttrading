package strategies

import (
	"math"

	at "agenttrading"
)

// S08TrendVolTarget: S01's trend filter, sized to target 40% annualised volatility using the
// 20-day realised volatility (capped at full size).
func S08TrendVolTarget(b at.Bars) []float64 {
	trend := S01TrendSMA(b)
	returns := make([]float64, b.Len())
	for i := 1; i < b.Len(); i++ {
		returns[i] = b.Close[i]/b.Close[i-1] - 1
	}
	volatility := at.StdDev(returns, 20)
	out := make([]float64, b.Len())
	for i := range out {
		if trend[i] > 0 && volatility[i] > 0 {
			out[i] = math.Min(1, 0.40/(volatility[i]*math.Sqrt(252)))
		}
	}
	return out
}
