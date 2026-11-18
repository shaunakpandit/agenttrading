package agenttrading

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

// Table formats results one per row, for comparing strategies at a glance.
func Table(results ...Result) string {
	var out strings.Builder
	w := tabwriter.NewWriter(&out, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Strategy\tReturn\tCAGR\tB&H\tMaxDD\tSharpe\tSortino\tTrades\tWin%\tPF\tAvgWin\tAvgLoss\tDays\tExposure\t")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%.2f\t%.2f\t%d\t%.0f\t%.2f\t%s\t%s\t%.1f\t%s\t\n",
			r.Name, percent(r.Return), percent(r.CAGR), percent(r.BuyHold), percent(-r.MaxDrawdown),
			r.Sharpe, r.Sortino, r.Trades, 100*r.WinRate, r.ProfitFactor,
			percent(r.AvgWin), percent(r.AvgLoss), r.AvgBars, percent(r.Exposure))
	}
	w.Flush()
	return out.String()
}

func percent(x float64) string { return fmt.Sprintf("%.1f%%", 100*x) }
