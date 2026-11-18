// Package strategies holds one Strategy per file. Register each in All.
package strategies

import at "agenttrading"

type Named struct {
	Name     string
	Strategy at.Strategy
}

var All = []Named{
	{"buy_and_hold", BuyAndHold},
}

// BuyAndHold is the baseline every attempt has to beat on risk-adjusted terms.
func BuyAndHold(b at.Bars) []float64 {
	out := make([]float64, b.Len())
	for i := range out {
		out[i] = 1
	}
	return out
}
