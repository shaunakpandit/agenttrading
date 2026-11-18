# First 10 attempts: TSLA daily swing (2026-11-18)

Data: Massive daily bars 2021-11-19 → 2026-11-17. In-sample 2022-09-01 → 2025-11-17, held-out
2025-11-18 → 2026-11-17 (looked at once, after all 10 were written). Cost 5 bps per side.

| Strategy | In-sample return / Sharpe / MaxDD | Held-out return / Sharpe |
|---|---|---|
| buy_and_hold | 47.5% / 0.50 / -65% | -5.3% / 0.11 |
| s01_trend_sma | 78.3% / 0.67 / -34% | -35.8% / -1.73 |
| s02_rsi2_pullback | 2.4% / 0.14 / -34% | -2.3% / -0.02 |
| s03_donchian (long/short) | 60.7% / 0.55 / -68% | -38.5% / -1.10 |
| s04_bollinger_reversion | -4.8% / 0.11 / -51% | **32.9% / 1.33** |
| s05_breakout_atr_trail | 96.5% / 0.76 / -40% | -12.0% / -0.33 |
| s06_macd_trend | **118.3% / 0.92 / -38%** | -38.0% / -1.83 |
| s07_squeeze_breakout | 12.1% / 0.30 / -18% | 0 trades |
| s08_trend_vol_target | 51.3% / 0.60 / -26% | -35.2% / -1.75 |
| s09_macd_atr_trail | 48.6% / 0.51 / -40% | -39.0% / -1.74 |
| s10_trend_ensemble | 105.9% / 0.87 / -31% | -29.0% / -1.43 |

## Findings

- Nothing here has a demonstrated edge. The in-sample winners (trend following) lost 30-39% in the
  held-out year, a choppy range where buy and hold lost 5%. The in-sample loser (s04 mean reversion)
  was the only winner. That's a regime change, and 10 attempts on 3 years of one stock can't tell
  skill from fit.
- On AMD, MU and NVDA (2022-09 → 2026-11, never used for development), no rule beats buy and hold
  on return in that bull market. s05, s06 and s09 come closest on Sharpe with much smaller drawdowns.
- s09 lesson: re-entering immediately after a stop while the entry condition is still true causes
  whipsaw clusters. Require a fresh trigger (a cross) or a cooldown.
- Fractional sizing (s08, s10) shows many "trades" because the trade log splits on every size change.

## Next ideas

Use a regime switch (trend rules when ADX or the 50-day slope is strong, s04-style reversion in
ranges), walk-forward over several tickers instead of one holdout, and get more history than the
5-year plan limit.
