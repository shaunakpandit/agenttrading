// Package agenttrading is a small toolkit for writing and testing daily-bar swing strategies.
// Fetch bars, write a Strategy, run Backtest and CheckLookahead, print Table. See AGENTS.md.
package agenttrading

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Bars is a daily OHLCV series, oldest first. Times are each session's date (midnight New York).
type Bars struct {
	Time                           []time.Time
	Open, High, Low, Close, Volume []float64
}

func (b Bars) Len() int { return len(b.Close) }

// Slice returns bars [i, j).
func (b Bars) Slice(i, j int) Bars {
	return Bars{b.Time[i:j], b.Open[i:j], b.High[i:j], b.Low[i:j], b.Close[i:j], b.Volume[i:j]}
}

const cacheDir = ".cache"

// Fetch returns split-adjusted daily bars for ticker between from and to (YYYY-MM-DD, inclusive)
// from Massive (formerly Polygon) using MASSIVE_API_KEY. Ranges that end before today are cached.
func Fetch(ticker, from, to string) (Bars, error) {
	path := filepath.Join(cacheDir, fmt.Sprintf("%s_%s_%s.json", ticker, from, to))
	if cached, err := os.ReadFile(path); err == nil {
		var b Bars
		return b, json.Unmarshal(cached, &b)
	}
	b, err := download(ticker, from, to)
	if err != nil {
		return Bars{}, err
	}
	if b.Len() == 0 {
		return b, fmt.Errorf("%s: no bars between %s and %s", ticker, from, to)
	}
	if to < time.Now().Format(time.DateOnly) { // today's bar is incomplete, so don't freeze it
		encoded, _ := json.Marshal(b)
		if err := os.MkdirAll(cacheDir, 0o755); err == nil {
			_ = os.WriteFile(path, encoded, 0o644) // cache failure only costs a refetch
		}
	}
	return b, nil
}

func download(ticker, from, to string) (Bars, error) {
	key := os.Getenv("MASSIVE_API_KEY")
	if key == "" {
		return Bars{}, errors.New("MASSIVE_API_KEY is not set")
	}
	url := fmt.Sprintf("https://api.massive.com/v2/aggs/ticker/%s/range/1/day/%s/%s?adjusted=true&sort=asc&limit=50000", ticker, from, to)
	newYork, _ := time.LoadLocation("America/New_York")
	var b Bars
	for url != "" {
		var page struct {
			Results []struct {
				T             int64
				O, H, L, C, V float64
			}
			NextURL string `json:"next_url"`
		}
		if err := getJSON(url, key, &page); err != nil {
			return Bars{}, fmt.Errorf("%s: %w", ticker, err)
		}
		for _, r := range page.Results {
			t := time.UnixMilli(r.T).In(newYork)
			b.Time = append(b.Time, time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, newYork))
			b.Open, b.High, b.Low = append(b.Open, r.O), append(b.High, r.H), append(b.Low, r.L)
			b.Close, b.Volume = append(b.Close, r.C), append(b.Volume, r.V)
		}
		url = page.NextURL
	}
	return b, nil
}

// getJSON waits out HTTP 429: the free plan allows 5 calls a minute.
func getJSON(url, key string, out any) error {
	for try := 1; ; try++ {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		request.Header.Set("Authorization", "Bearer "+key)
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return err
		}
		body, _ := io.ReadAll(response.Body)
		response.Body.Close()
		if response.StatusCode == http.StatusTooManyRequests && try < 12 {
			fmt.Fprintln(os.Stderr, "note: rate limited, retrying in 15s")
			time.Sleep(15 * time.Second)
			continue
		}
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("%s: %.300s", response.Status, body)
		}
		return json.Unmarshal(body, out)
	}
}

// Stats summarises the last year of b: average daily dollar volume and annualised volatility
// of daily returns. Use it to pick a liquid, volatile ticker.
func Stats(b Bars) (avgDollarVolume, annualVolatility float64) {
	recent := b.Slice(max(b.Len()-252, 0), b.Len())
	returns := make([]float64, 0, recent.Len())
	for i := range recent.Close {
		avgDollarVolume += recent.Close[i] * recent.Volume[i] / float64(recent.Len())
		if i > 0 {
			returns = append(returns, recent.Close[i]/recent.Close[i-1]-1)
		}
	}
	_, deviation, _ := moments(returns)
	return avgDollarVolume, deviation * math.Sqrt(252)
}
