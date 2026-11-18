// Command run backtests every registered strategy on one ticker.
//
//	go run ./cmd/run -ticker TSLA                 in-sample only (develop here)
//	go run ./cmd/run -ticker TSLA -oos            also the held-out window
//	go run ./cmd/run -ticker TSLA -trades s03_x   print one strategy's trades
//	go run ./cmd/run -screen TSLA,NVDA,AMD        dollar volume and volatility per ticker
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	at "agenttrading"
	"agenttrading/strategies"
)

func main() {
	yesterday := time.Now().AddDate(0, 0, -1).Format(time.DateOnly)
	ticker := flag.String("ticker", "TSLA", "ticker to backtest")
	from := flag.String("from", "2016-01-01", "first date fetched (includes warm-up)")
	to := flag.String("to", yesterday, "last date fetched")
	start := flag.String("start", "2017-01-01", "first scored date; earlier bars are warm-up")
	holdout := flag.String("holdout", "2024-11-18", "first out-of-sample date; in-sample ends the day before")
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

	b, err := at.Fetch(*ticker, *from, *to)
	fail(err)
	startDate, holdoutDate := date(*start), date(*holdout)
	inSample := at.Config{CostBps: *cost, From: startDate, To: holdoutDate.AddDate(0, 0, -1)}
	outOfSample := at.Config{CostBps: *cost, From: holdoutDate}

	var in, out []at.Result
	for _, s := range strategies.All {
		if err := at.CheckLookahead(b, s.Strategy); err != nil {
			fail(fmt.Errorf("%s: %w", s.Name, err))
		}
		in = append(in, at.Backtest(s.Name, b, s.Strategy, inSample))
		out = append(out, at.Backtest(s.Name, b, s.Strategy, outOfSample))
		if s.Name == *trades {
			for _, t := range in[len(in)-1].TradeLog {
				fmt.Printf("%+.0f  %s %8.2f → %s %8.2f  %6.1f%%  %3d days\n", t.Side,
					t.EntryTime.Format(time.DateOnly), t.EntryPrice, t.ExitTime.Format(time.DateOnly), t.ExitPrice, 100*t.Return, t.Bars)
			}
		}
	}
	fmt.Printf("%s in-sample %s → %s, cost %.0f bps/side\n%s", *ticker, *start, inSample.To.Format(time.DateOnly), *cost, at.Table(in...))
	if *oos {
		fmt.Printf("\n%s out-of-sample %s → %s\n%s", *ticker, *holdout, *to, at.Table(out...))
	}
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
