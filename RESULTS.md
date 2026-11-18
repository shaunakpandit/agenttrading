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

# Round 2: 10 more attempts across 10 stocks (2026-11-18)

Universe: TSLA, NVDA, AMD, MU, META, AVGO, PLTR, COIN, MSTR, SMCI, all liquid and volatile. Same
split: in-sample 2022-09-01 → 2025-11-17 (a bull market, buy and hold median Sharpe 1.16) and a
held-out year 2025-11-18 → 2026-11-17 (choppy, buy and hold median Sharpe 0.14), looked at once.
Ranked by median Sharpe across tickers, not by one stock.

| Strategy | In-sample med Sharpe / MaxDD / beat B&H | Held-out med Sharpe / MaxDD / beat B&H |
|---|---|---|
| buy_and_hold | 1.16 / -58% / – | 0.14 / -42% / – |
| s05_breakout_atr_trail (round 1) | 1.09 / -40% / 5 | -0.61 / -35% / 1 |
| s11_regime_adx | 0.99 / -45% / 1 | -0.14 / -30% / 2 |
| s12_fresh_cross_trend | **1.22** / -36% / 4 | -0.71 / -27% / 2 |
| s13_trend_pullback | 0.23 / -36% / 1 | -0.87 / -15% / 1 |
| s14_ibs_reversion | 0.71 / -29% / 2 | **0.31 / -15% / 8** |
| s15_time_series_momentum | 0.99 / -45% / 3 | -0.40 / -29% / 1 |
| s16_volume_breakout | 0.92 / -39% / 1 | -0.36 / -26% / 2 |
| s17_efficient_breakout | 0.98 / -31% / 3 | -0.28 / -20% / 2 |
| s18_breakout_long_trend | 1.11 / -35% / 5 | -0.42 / -29% / 2 |
| s19_breakout_cross_ensemble | 1.18 / -34% / 5 | -0.70 / -27% / 2 |
| s20_trend_reversion_ensemble | 1.11 / -30% / 3 | -0.28 / -19% / 2 |

## Findings

- Trend rules are regime bets. All of them nearly matched buy and hold's Sharpe in the bull years
  with lower drawdown, then all lost in the choppy year. Filters for volume, efficiency, a long-term
  trend or ADX didn't change that.
- s14 (IBS < 0.15 above SMA(100), exit above the prior high or after 5 days) is the only rule that
  was positive in both regimes on most names. It's modest: 14-27% exposure, about 3-day holds,
  ~70% win rate.
- The stock that could work is **AMD with s14**: Sharpe 1.10 in-sample (-15% DD) and 1.26 held-out
  (+52.5%, -18% DD, 23 trades). Caveats: AMD was chosen after seeing the held-out year, and one
  earnings-gap trade (+23% in May 2026) supplies roughly half of the held-out return. A clean test
  needs data that doesn't exist yet: paper trade it forward.

## Next ideas

Tune nothing further on this data. Paper trade s14 on AMD (and the universe) forward. Try IBS
variants built on other evidence: earnings-date avoidance, and a market (SPY/QQQ) regime filter,
which needs strategies to see a second series.
