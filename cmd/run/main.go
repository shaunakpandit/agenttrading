// Command run backtests every registered strategy on one ticker.
//
//	go run ./cmd/run -ticker TSLA                 in-sample only (develop here)
//	go run ./cmd/run -ticker TSLA,NVDA,AMD        per-strategy summary and Sharpe per ticker
//	go run ./cmd/run -ticker TSLA -oos            also the held-out window
//	go run ./cmd/run -ticker TSLA -trades s03_x   print one strategy's trades
//	go run ./cmd/run -screen TSLA,NVDA,AMD        dollar volume and volatility per ticker
package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	at "agenttrading"
	"agenttrading/strategies"
)

func main() {
	yesterday := time.Now().AddDate(0, 0, -1).Format(time.DateOnly)
	tickerList := flag.String("ticker", "TSLA", "ticker, or comma-separated tickers, to backtest")
	from := flag.String("from", "2016-01-01", "first date fetched (includes warm-up)")
	to := flag.String("to", yesterday, "last date fetched")
	start := flag.String("start", "2022-09-01", "first scored date; earlier bars are warm-up")
	holdout := flag.String("holdout", "2025-11-18", "first out-of-sample date; in-sample ends the day before")
	oos := flag.Bool("oos", false, "also score the out-of-sample window")
	cost := flag.Float64("cost", 5, "cost per side in basis points")
	trades := flag.String("trades", "", "print the in-sample trades of this strategy")
	screen := flag.String("screen", "", "comma-separated tickers to screen instead of backtesting")
	flag.Parse()

	if *screen != "" {
		for _, t := range strings.Split(*screen, ",") {
			b, err := at.Fetch(t, *from, *to)
			fail(err)
			volume, volatility := at.Stats(b)
			fmt.Printf("%-6s  $%6.2fB/day  vol %5.1f%%  last %.2f\n", t, volume/1e9, 100*volatility, b.Close[b.Len()-1])
		}
		return
	}

	startDate, holdoutDate := date(*start), date(*holdout)
	inSample := at.Config{CostBps: *cost, From: startDate, To: holdoutDate.AddDate(0, 0, -1)}
	outOfSample := at.Config{CostBps: *cost, From: holdoutDate}
	tickers := strings.Split(*tickerList, ",")

	// in[s][t] and out[s][t] are strategy s on ticker t; buyHold[t] is the per-ticker baseline.
	in, out := make([][]at.Result, len(strategies.All)), make([][]at.Result, len(strategies.All))
	var buyHoldIn, buyHoldOut []at.Result
	for _, ticker := range tickers {
		b, err := at.Fetch(ticker, *from, *to)
		fail(err)
		buyHoldIn = append(buyHoldIn, at.Backtest("buy_and_hold", b, strategies.BuyAndHold, inSample))
		buyHoldOut = append(buyHoldOut, at.Backtest("buy_and_hold", b, strategies.BuyAndHold, outOfSample))
		for i, s := range strategies.All {
			if err := at.CheckLookahead(b, s.Strategy); err != nil {
				fail(fmt.Errorf("%s on %s: %w", s.Name, ticker, err))
			}
			in[i] = append(in[i], at.Backtest(s.Name, b, s.Strategy, inSample))
			out[i] = append(out[i], at.Backtest(s.Name, b, s.Strategy, outOfSample))
			if s.Name == *trades {
				fmt.Println(ticker)
				for _, t := range in[i][len(in[i])-1].TradeLog {
					fmt.Printf("%+.0f  %s %8.2f → %s %8.2f  %6.1f%%  %3d days\n", t.Side,
						t.EntryTime.Format(time.DateOnly), t.EntryPrice, t.ExitTime.Format(time.DateOnly), t.ExitPrice, 100*t.Return, t.Bars)
				}
			}
		}
	}
	report := func(title string, results [][]at.Result, buyHold []at.Result) {
		fmt.Printf("%s %s, cost %.0f bps/side\n", *tickerList, title, *cost)
		if len(tickers) == 1 {
			var rows []at.Result
			for _, r := range results {
				rows = append(rows, r[0])
			}
			fmt.Print(at.Table(rows...))
			return
		}
		fmt.Print(summary(tickers, results, buyHold))
	}
	report(fmt.Sprintf("in-sample %s → %s", *start, inSample.To.Format(time.DateOnly)), in, buyHoldIn)
	if *oos {
		fmt.Println()
		report(fmt.Sprintf("out-of-sample %s → %s", *holdout, *to), out, buyHoldOut)
	}
}

// summary prints one row per strategy: medians across tickers, how many tickers it beat buy and
// hold's Sharpe on, and its Sharpe on each ticker.
func summary(tickers []string, results [][]at.Result, buyHold []at.Result) string {
	var text strings.Builder
	w := tabwriter.NewWriter(&text, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Strategy\tMedSharpe\tMedReturn\tMedMaxDD\tMedExposure\tBeatB&H\t%s\t\n", strings.Join(tickers, "\t"))
	row := func(name string, rs []at.Result, beat int) {
		cells := make([]string, len(rs))
		for i, r := range rs {
			cells[i] = fmt.Sprintf("%.2f", r.Sharpe)
		}
		fmt.Fprintf(w, "%s\t%.2f\t%.1f%%\t%.1f%%\t%.0f%%\t%d/%d\t%s\t\n", name,
			median(rs, func(r at.Result) float64 { return r.Sharpe }),
			100*median(rs, func(r at.Result) float64 { return r.Return }),
			-100*median(rs, func(r at.Result) float64 { return r.MaxDrawdown }),
			100*median(rs, func(r at.Result) float64 { return r.Exposure }),
			beat, len(rs), strings.Join(cells, "\t"))
	}
	row("buy_and_hold", buyHold, 0)
	for _, rs := range results {
		beat := 0
		for i, r := range rs {
			if r.Sharpe > buyHold[i].Sharpe {
				beat++
			}
		}
		row(rs[0].Name, rs, beat)
	}
	w.Flush()
	return text.String()
}

func median(results []at.Result, field func(at.Result) float64) float64 {
	values := make([]float64, len(results))
	for i, r := range results {
		values[i] = field(r)
	}
	slices.Sort(values)
	if n := len(values); n%2 == 0 {
		return (values[n/2-1] + values[n/2]) / 2
	}
	return values[len(values)/2]
}

func date(s string) time.Time {
	newYork, err := time.LoadLocation("America/New_York") // bar times are New York midnights
	fail(err)
	t, err := time.ParseInLocation(time.DateOnly, s, newYork)
	fail(err)
	return t
}

func fail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
