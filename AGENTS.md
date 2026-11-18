# agenttrading

A stdlib-only Go toolkit for writing and testing daily-bar swing strategies on US stocks.

## Setup

`MASSIVE_API_KEY` must be set (Massive, formerly Polygon). You can also put it in `.env` and run
`set -a; . ./.env; set +a`. Bars are cached in `.cache/`, so repeat runs make no network calls.

## Loop for a new strategy

1. Add `strategies/sNN_name.go` with `func Name(b at.Bars) []float64` and one line in `strategies.All`.
2. Run `go run ./cmd/run -ticker TSLA`. Every strategy is lookahead-checked, then scored in-sample.
3. Run `go run ./cmd/run -ticker TSLA -trades sNN_name` to inspect its trades.
4. Look at the held-out window (`-oos`) only once you're done iterating. Each peek is another fit.

Pick a ticker with `go run ./cmd/run -screen TSLA,NVDA,AMD,PLTR` (dollar volume, volatility).

## Contract

- A strategy returns one value per bar, the target position after bar i's close: 1 long, 0 flat,
  -1 short, fractions for partial size. NaN means 0, and values are clamped to [-1, 1].
- Value i may only use bars 0..i. `CheckLookahead` enforces this by truncating the data.
- Fills happen at bar i's close, and the position earns bar i+1's close-to-close move. Stops act on
  closes, never intrabar. Each change in position costs `CostBps` per unit (5 bps per side by default).

## API (package `agenttrading`, import as `at`)

- Data: `Fetch(ticker, from, to)`, `Bars{Time, Open, High, Low, Close, Volume}`, `b.Len()`,
  `b.Slice(i, j)`, `Stats(b)`.
- Indicators (same length as the input, NaN during warm-up): `SMA EMA RMA StdDev Highest Lowest
  Shift RSI ATR MACD Bollinger CrossOver CrossUnder`.
- Rules: `EnterExit(enter, exit []bool, side)` holds `side` from an enter until the next exit.
  Use a plain loop when a rule needs state, such as a trailing stop.
- Testing: `Backtest(name, b, s, Config{CostBps, From, To})` returns a `Result` with `Return`,
  `CAGR`, `BuyHold`, `MaxDrawdown`, `Sharpe`, `Sortino`, `Trades`, `WinRate`, `ProfitFactor`,
  `AvgWin`, `AvgLoss`, `AvgBars`, `Exposure`, `Equity` and `TradeLog`. `Table(results...)` formats them.
  `CheckLookahead(b, s)` returns an error if a signal looks ahead.

Run `go test ./...` for tests on synthetic data (no network).
