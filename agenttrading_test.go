package agenttrading

import (
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

func randomWalk(n int) Bars {
	rng := rand.New(rand.NewPCG(7, 8))
	var b Bars
	price, day := 100.0, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range n {
		open := price
		price *= 1 + rng.NormFloat64()*0.02
		b.Time = append(b.Time, day.AddDate(0, 0, i))
		b.Open, b.Close = append(b.Open, open), append(b.Close, price)
		b.High, b.Low = append(b.High, max(open, price)*1.01), append(b.Low, min(open, price)*0.99)
		b.Volume = append(b.Volume, 1e6)
	}
	return b
}

func always(v float64) Strategy {
	return func(b Bars) []float64 {
		out := make([]float64, b.Len())
		for i := range out {
			out[i] = v
		}
		return out
	}
}

func TestAlwaysLongMatchesBuyHoldMinusOneFee(t *testing.T) {
	b := randomWalk(500)
	r := Backtest("long", b, always(1), Config{CostBps: 10})
	// the fee comes off the first bar's return, then compounds
	want := (b.Close[1]/b.Close[0]-0.001)*b.Close[499]/b.Close[1] - 1
	if math.Abs(r.Return-want) > 1e-9 || r.Trades != 1 {
		t.Fatalf("return %v trades %d, want %v and 1", r.Return, r.Trades, want)
	}
}

func TestWindowScoresOnlyItsBars(t *testing.T) {
	b := randomWalk(500)
	r := Backtest("long", b, always(1), Config{From: b.Time[100], To: b.Time[199]})
	if len(r.Equity) != 100 || math.Abs(r.BuyHold-(b.Close[199]/b.Close[100]-1)) > 1e-12 {
		t.Fatalf("equity len %d buyhold %v", len(r.Equity), r.BuyHold)
	}
}

func TestCheckLookahead(t *testing.T) {
	b := randomWalk(500)
	peek := func(b Bars) []float64 {
		out := make([]float64, b.Len())
		for i := 0; i+1 < b.Len(); i++ {
			if b.Close[i+1] > b.Close[i] {
				out[i] = 1
			}
		}
		return out
	}
	if CheckLookahead(b, peek) == nil {
		t.Fatal("peeking strategy passed")
	}
	fair := func(b Bars) []float64 {
		return EnterExit(CrossOver(b.Close, SMA(b.Close, 20)), CrossUnder(b.Close, SMA(b.Close, 20)), 1)
	}
	if err := CheckLookahead(b, fair); err != nil {
		t.Fatal(err)
	}
}

func TestEnterExit(t *testing.T) {
	got := EnterExit([]bool{false, true, false, true, false}, []bool{true, false, true, true, false}, -1)
	want := []float64{0, -1, 0, -1, -1}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}

func TestIndicators(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	if s := SMA(x, 3); !math.IsNaN(s[1]) || s[4] != 4 {
		t.Fatalf("SMA %v", s)
	}
	if s := Shift(x, 2); !math.IsNaN(s[1]) || s[4] != 3 {
		t.Fatalf("Shift %v", s)
	}
	if r := RSI(x, 2); r[4] != 100 {
		t.Fatalf("RSI %v", r)
	}
	// a steady uptrend has no down moves, so ADX converges to 100
	b := randomWalk(60)
	for i := range b.Close {
		b.Close[i], b.High[i], b.Low[i] = float64(100+i), float64(101+i), float64(99+i)
	}
	if a := ADX(b, 14); !math.IsNaN(a[10]) || math.Abs(a[59]-100) > 1e-9 {
		t.Fatalf("ADX %v", a)
	}
}
