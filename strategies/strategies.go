// Package strategies holds one Strategy per file. Register each in All. cmd/run always adds the
// BuyAndHold baseline itself.
package strategies

import at "agenttrading"

type Named struct {
	Name     string
	Strategy at.Strategy
}

var All = []Named{
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
	{"s11_regime_adx", S11RegimeADX},
	{"s12_fresh_cross_trend", S12FreshCrossTrend},
	{"s13_trend_pullback", S13TrendPullback},
	{"s14_ibs_reversion", S14IBSReversion},
	{"s15_time_series_momentum", S15TimeSeriesMomentum},
	{"s16_volume_breakout", S16VolumeBreakout},
	{"s17_efficient_breakout", S17EfficientBreakout},
	{"s18_breakout_long_trend", S18BreakoutLongTrend},
	{"s19_breakout_cross_ensemble", S19BreakoutCrossEnsemble},
	{"s20_trend_reversion_ensemble", S20TrendReversionEnsemble},
}

// BuyAndHold is the baseline every attempt has to beat on risk-adjusted terms.
func BuyAndHold(b at.Bars) []float64 {
	out := make([]float64, b.Len())
	for i := range out {
		out[i] = 1
	}
	return out
}
