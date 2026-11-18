package agenttrading

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

// Strategy returns a target position per bar: 1 long, 0 flat, -1 short, fractions in between
// for partial size, NaN treated as 0. Value i may use only bars 0..i.
type Strategy func(b Bars) []float64

type Config struct {
	CostBps  float64   // per side, fee plus slippage, in basis points
	From, To time.Time // score only this window (zero = unbounded); earlier bars are warm-up
}

type Trade struct {
	Side                  float64
	EntryTime, ExitTime   time.Time
	EntryPrice, ExitPrice float64
	Return                float64 // net of entry and exit costs
	Bars                  int
}

// Result fields are fractions (0.05 = 5%) unless noted.
type Result struct {
	Name                                   string
	Return, CAGR, BuyHold, MaxDrawdown     float64
	Sharpe, Sortino                        float64
	Trades                                 int
	WinRate, ProfitFactor, AvgWin, AvgLoss float64
	AvgBars, Exposure                      float64 // AvgBars in days held per trade
	Times                                  []time.Time
	Equity                                 []float64 // starts at 1
	TradeLog                               []Trade
}

// Backtest runs s on all of b and scores the window [cfg.From, cfg.To]. The signal at bar i's
// close is filled at that close and earns bar i+1's close-to-close move; each change in position
// pays cfg.CostBps per unit. The position is flat at the window start.
func Backtest(name string, b Bars, s Strategy, cfg Config) Result {
	signal := s(b)
	if len(signal) != b.Len() {
		panic(fmt.Sprintf("strategy %s returned %d signals for %d bars", name, len(signal), b.Len()))
	}
	start, end := 0, b.Len()
	for start < end && b.Time[start].Before(cfg.From) {
		start++
	}
	for !cfg.To.IsZero() && end > start && b.Time[end-1].After(cfg.To) {
		end--
	}
	if end-start < 2 {
		return Result{Name: name}
	}
	closes, times := b.Close[start:end], b.Time[start:end]
	cost := cfg.CostBps / 1e4
	positions := make([]float64, len(closes))
	returns := make([]float64, len(closes)-1)
	equity := make([]float64, len(closes))
	equity[0] = 1
	for i := 1; i < len(closes); i++ {
		positions[i] = position(signal[start+i-1])
		returns[i-1] = positions[i]*(closes[i]/closes[i-1]-1) - math.Abs(positions[i]-positions[i-1])*cost
		equity[i] = equity[i-1] * (1 + returns[i-1])
	}
	r := Result{Name: name, Times: times, Equity: equity, TradeLog: tradeLog(positions, times, closes, cost)}
	score(&r, returns, positions, closes)
	return r
}

func position(signal float64) float64 {
	if math.IsNaN(signal) {
		return 0
	}
	return math.Max(-1, math.Min(1, signal))
}

// tradeLog splits positions into runs of equal non-zero size. A trade held from bar i was filled
// at the close of bar i-1; an open trade is closed at the last bar.
func tradeLog(positions []float64, times []time.Time, closes []float64, cost float64) []Trade {
	var trades []Trade
	entry := -1
	for i := 1; i <= len(positions); i++ {
		current := 0.0
		if i < len(positions) {
			current = positions[i]
		}
		if current == positions[i-1] {
			continue
		}
		if entry >= 0 {
			side := positions[i-1]
			trades = append(trades, Trade{Side: side, EntryTime: times[entry-1], EntryPrice: closes[entry-1],
				ExitTime: times[i-1], ExitPrice: closes[i-1], Bars: i - entry,
				Return: side*(closes[i-1]/closes[entry-1]-1) - 2*math.Abs(side)*cost})
			entry = -1
		}
		if current != 0 {
			entry = i
		}
	}
	return trades
}

func score(r *Result, returns, positions, closes []float64) {
	final := r.Equity[len(r.Equity)-1]
	r.Return, r.BuyHold = final-1, closes[len(closes)-1]/closes[0]-1
	if years := float64(len(returns)) / 252; final > 0 {
		r.CAGR = math.Pow(final, 1/years) - 1
	}
	peak := 1.0
	for _, e := range r.Equity {
		peak = max(peak, e)
		r.MaxDrawdown = max(r.MaxDrawdown, 1-e/peak)
	}
	m, deviation, downside := moments(returns)
	if deviation > 0 {
		r.Sharpe = m / deviation * math.Sqrt(252)
	}
	if downside > 0 {
		r.Sortino = m / downside * math.Sqrt(252)
	}
	for _, p := range positions {
		if p != 0 {
			r.Exposure += 1 / float64(len(positions))
		}
	}
	r.Trades = len(r.TradeLog)
	var wins, gains, losses, pains, bars float64
	for _, t := range r.TradeLog {
		bars += float64(t.Bars)
		if t.Return > 0 {
			wins, gains = wins+1, gains+t.Return
		} else {
			losses, pains = losses+1, pains-t.Return
		}
	}
	if r.Trades > 0 {
		r.WinRate, r.AvgBars = wins/float64(r.Trades), bars/float64(r.Trades)
	}
	r.AvgWin, r.AvgLoss, r.ProfitFactor = safeDivide(gains, wins), -safeDivide(pains, losses), safeDivide(gains, pains)
	if pains == 0 && gains > 0 {
		r.ProfitFactor = math.Inf(1)
	}
}

// moments returns the mean, population standard deviation and downside deviation of x.
func moments(x []float64) (m, deviation, downside float64) {
	n := float64(len(x))
	for _, v := range x {
		m += v / n
	}
	for _, v := range x {
		deviation += (v - m) * (v - m) / n
		downside += math.Pow(min(v, 0), 2) / n
	}
	return m, math.Sqrt(deviation), math.Sqrt(downside)
}

func safeDivide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

// CheckLookahead runs s on b, then on b truncated at five random points, and fails if any
// signal at or before the cut changes: signal i must not depend on bars after i. It also checks
// the output length and that values lie in [-1, 1] or are NaN.
func CheckLookahead(b Bars, s Strategy) error {
	full := s(b)
	if len(full) != b.Len() {
		return fmt.Errorf("returned %d signals for %d bars", len(full), b.Len())
	}
	for i, v := range full {
		if v < -1 || v > 1 {
			return fmt.Errorf("signal[%d] = %v, want a value in [-1, 1] or NaN", i, v)
		}
	}
	rng := rand.New(rand.NewPCG(1, 2))
	for range 5 {
		k := b.Len()/4 + rng.IntN(b.Len()*3/4)
		cut := s(b.Slice(0, k))
		for i := range cut {
			if !(cut[i] == full[i] || math.IsNaN(cut[i]) && math.IsNaN(full[i])) {
				return fmt.Errorf("lookahead: signal[%d] is %v on all bars but %v on the first %d bars", i, full[i], cut[i], k)
			}
		}
	}
	return nil
}
