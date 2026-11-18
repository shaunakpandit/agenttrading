package agenttrading

// Indicators. Every function returns a slice of the input's length with NaN where the value is
// not yet defined, and value i depends only on inputs 0..i.

import "math"

func nans(n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = math.NaN()
	}
	return out
}

func window(x []float64, n int, f func([]float64) float64) []float64 {
	out := nans(len(x))
	for i := n - 1; n > 0 && i < len(x); i++ {
		out[i] = f(x[i-n+1 : i+1])
	}
	return out
}

func zip(a, b []float64, f func(a, b float64) float64) []float64 {
	out := make([]float64, len(a))
	for i := range out {
		out[i] = f(a[i], b[i])
	}
	return out
}

func apply(x []float64, f func(float64) float64) []float64 {
	out := make([]float64, len(x))
	for i, v := range x {
		out[i] = f(v)
	}
	return out
}

func mean(x []float64) float64 {
	m, _, _ := moments(x)
	return m
}

func stdDev(x []float64) float64 {
	_, d, _ := moments(x)
	return d
}

func reduce(x []float64, f func(a, b float64) float64) float64 {
	out := x[0]
	for _, v := range x[1:] {
		out = f(out, v)
	}
	return out
}

// smooth is an exponential average seeded with the SMA of the first n valid values; NaN inputs give NaN.
func smooth(x []float64, n int, alpha float64) []float64 {
	out := nans(len(x))
	value, count := 0.0, 0
	for i, v := range x {
		switch {
		case math.IsNaN(v):
		case count < n:
			value, count = value+v/float64(n), count+1
			if count == n {
				out[i] = value
			}
		default:
			value = alpha*v + (1-alpha)*value
			out[i] = value
		}
	}
	return out
}

func SMA(x []float64, n int) []float64    { return window(x, n, mean) }
func StdDev(x []float64, n int) []float64 { return window(x, n, stdDev) }
func Highest(x []float64, n int) []float64 {
	return window(x, n, func(w []float64) float64 { return reduce(w, math.Max) })
}
func Lowest(x []float64, n int) []float64 {
	return window(x, n, func(w []float64) float64 { return reduce(w, math.Min) })
}
func EMA(x []float64, n int) []float64 { return smooth(x, n, 2/float64(n+1)) }
func RMA(x []float64, n int) []float64 { return smooth(x, n, 1/float64(n)) }

// Shift returns x delayed by n bars (out[i] = x[i-n]), e.g. Shift(Highest(high, 20), 1) for a
// breakout level that excludes today.
func Shift(x []float64, n int) []float64 {
	out := nans(len(x))
	for i := n; i < len(x); i++ {
		out[i] = x[i-n]
	}
	return out
}

// RSI uses Wilder smoothing.
func RSI(x []float64, n int) []float64 {
	change := zip(x, Shift(x, 1), func(a, b float64) float64 { return a - b })
	gain := RMA(apply(change, func(c float64) float64 { return math.Max(c, 0) }), n)
	loss := RMA(apply(change, func(c float64) float64 { return math.Max(-c, 0) }), n)
	return zip(gain, loss, func(g, l float64) float64 {
		if l == 0 {
			return 100
		}
		return 100 - 100/(1+g/l)
	})
}

// ATR uses Wilder smoothing of the true range.
func ATR(b Bars, n int) []float64 {
	tr := make([]float64, b.Len())
	for i := range tr {
		tr[i] = b.High[i] - b.Low[i]
		if i > 0 {
			tr[i] = max(tr[i], math.Abs(b.High[i]-b.Close[i-1]), math.Abs(b.Low[i]-b.Close[i-1]))
		}
	}
	return RMA(tr, n)
}

func MACD(x []float64, fast, slow, signal int) (macd, sig, hist []float64) {
	macd = zip(EMA(x, fast), EMA(x, slow), func(a, b float64) float64 { return a - b })
	sig = EMA(macd, signal)
	return macd, sig, zip(macd, sig, func(a, b float64) float64 { return a - b })
}

func Bollinger(x []float64, n int, k float64) (mid, up, lo []float64) {
	mid, sd := SMA(x, n), StdDev(x, n)
	return mid, zip(mid, sd, func(m, s float64) float64 { return m + k*s }), zip(mid, sd, func(m, s float64) float64 { return m - k*s })
}

func cross(a, b []float64, above bool) []bool {
	out := make([]bool, len(a))
	for i := 1; i < len(a); i++ {
		if above {
			out[i] = a[i] > b[i] && a[i-1] <= b[i-1]
		} else {
			out[i] = a[i] < b[i] && a[i-1] >= b[i-1]
		}
	}
	return out
}

func CrossOver(a, b []float64) []bool  { return cross(a, b, true) }
func CrossUnder(a, b []float64) []bool { return cross(a, b, false) }

// EnterExit turns "enter when X, exit when Y" rules into positions: side (1 long, -1 short) from
// the first enter until the next exit. Exit wins if both fire on the same bar while in a position.
func EnterExit(enter, exit []bool, side float64) []float64 {
	out := make([]float64, len(enter))
	held := 0.0
	for i := range enter {
		switch {
		case held != 0 && exit[i]:
			held = 0
		case held == 0 && enter[i]:
			held = side
		}
		out[i] = held
	}
	return out
}

// ADX is Wilder's average directional index: trend strength from 0 to 100, regardless of direction.
// Above about 25 the market is trending; below about 20 it is ranging.
func ADX(b Bars, n int) []float64 {
	plus, minus := make([]float64, b.Len()), make([]float64, b.Len())
	plus[0], minus[0] = math.NaN(), math.NaN()
	for i := 1; i < b.Len(); i++ {
		up, down := b.High[i]-b.High[i-1], b.Low[i-1]-b.Low[i]
		if up > down && up > 0 {
			plus[i] = up
		}
		if down > up && down > 0 {
			minus[i] = down
		}
	}
	atr, plusSmoothed, minusSmoothed := ATR(b, n), RMA(plus, n), RMA(minus, n)
	dx := make([]float64, b.Len())
	for i := range dx {
		plusDI, minusDI := plusSmoothed[i]/atr[i], minusSmoothed[i]/atr[i]
		dx[i] = 100 * math.Abs(plusDI-minusDI) / (plusDI + minusDI)
	}
	return RMA(dx, n)
}
