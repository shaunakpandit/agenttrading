// Package strategies holds one Strategy per file. Register each in All.
package strategies

import at "agenttrading"

type Named struct {
	Name     string
	Strategy at.Strategy
}

var All = []Named{
	{"buy_and_hold", BuyAndHold},
	{"s01_trend_sma", S01TrendSMA},
	{"s02_rsi2_pullback", S02RSI2Pullback},
	{"s03_donchian", S03Donchian},
	{"s04_bollinger_reversion", S04BollingerReversion},
	{"s05_breakout_atr_trail", S05BreakoutATRTrail},
	{"s06_macd_trend", S06MACDTrend},
	{"s07_squeeze_breakout", S07SqueezeBreakout},
	{"s08_trend_vol_target", S08TrendVolTarget},
	{"s09_macd_atr_trail", S09MACDATRTrail},
	{"s10_trend_ensemble", S10TrendEnsemble},
}

// BuyAndHold is the baseline every attempt has to beat on risk-adjusted terms.
func BuyAndHold(b at.Bars) []float64 {
	out := make([]float64, b.Len())
	for i := range out {
		out[i] = 1
	}
	return out
}
