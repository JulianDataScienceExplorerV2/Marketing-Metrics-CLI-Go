// Marketing Metrics CLI
// ======================
// Reads a CSV file containing campaign performance data and computes
// standard marketing KPIs: CTR, CVR, CPC, CPA, ROAS, and Revenue/Conv.
//
// Usage:
//
//	go run main.go --file data/sample_campaigns.csv
//	go run main.go --file data/sample_campaigns.csv --sort roas
//	go run main.go --file data/sample_campaigns.csv --min-roas 2.0
//
// CSV expected columns (header row required):
//
//	campaign_name, impressions, clicks, conversions, spend, revenue
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

// ── Types ──────────────────────────────────────────────────────────────────

// Campaign holds raw data parsed from one CSV row.
type Campaign struct {
	Name        string
	Impressions float64
	Clicks      float64
	Conversions float64
	Spend       float64
	Revenue     float64
}

// Metrics holds computed KPIs for one campaign.
type Metrics struct {
	Name        string
	CTR         float64 // Click-Through Rate        = clicks / impressions * 100
	CVR         float64 // Conversion Rate            = conversions / clicks * 100
	CPC         float64 // Cost Per Click             = spend / clicks
	CPA         float64 // Cost Per Acquisition       = spend / conversions
	ROAS        float64 // Return on Ad Spend         = revenue / spend
	RevPerConv  float64 // Revenue Per Conversion     = revenue / conversions
	Spend       float64
	Revenue     float64
	Impressions float64
	Clicks      float64
	Conversions float64
}

// ── Helpers ────────────────────────────────────────────────────────────────

func safe(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func computeMetrics(c Campaign) Metrics {
	return Metrics{
		Name:        c.Name,
		CTR:         round2(safe(c.Clicks, c.Impressions) * 100),
		CVR:         round2(safe(c.Conversions, c.Clicks) * 100),
		CPC:         round2(safe(c.Spend, c.Clicks)),
		CPA:         round2(safe(c.Spend, c.Conversions)),
		ROAS:        round2(safe(c.Revenue, c.Spend)),
		RevPerConv:  round2(safe(c.Revenue, c.Conversions)),
		Spend:       c.Spend,
		Revenue:     c.Revenue,
		Impressions: c.Impressions,
		Clicks:      c.Clicks,
		Conversions: c.Conversions,
	}
}

// ── CSV Parsing ────────────────────────────────────────────────────────────

var requiredCols = []string{"campaign_name", "impressions", "clicks", "conversions", "spend", "revenue"}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", ""))
	return strconv.ParseFloat(s, 64)
}

func loadCSV(path string) ([]Campaign, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSV parse error: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("CSV must have a header row and at least one data row")
	}

	// Build column index from header
	header := make(map[string]int)
	for i, h := range rows[0] {
		header[strings.ToLower(strings.TrimSpace(h))] = i
	}
	for _, col := range requiredCols {
		if _, ok := header[col]; !ok {
			return nil, fmt.Errorf("missing required column: %q", col)
		}
	}

	var campaigns []Campaign
	for i, row := range rows[1:] {
		if len(row) == 0 {
			continue
		}
		parse := func(col string) (float64, error) {
			v, err := parseFloat(row[header[col]])
			if err != nil {
				return 0, fmt.Errorf("row %d, column %q: %w", i+2, col, err)
			}
			return v, nil
		}

		imp, err := parse("impressions")
		if err != nil {
			return nil, err
		}
		clk, err := parse("clicks")
		if err != nil {
			return nil, err
		}
		conv, err := parse("conversions")
		if err != nil {
			return nil, err
		}
		spend, err := parse("spend")
		if err != nil {
			return nil, err
		}
		rev, err := parse("revenue")
		if err != nil {
			return nil, err
		}

		campaigns = append(campaigns, Campaign{
			Name:        strings.TrimSpace(row[header["campaign_name"]]),
			Impressions: imp,
			Clicks:      clk,
			Conversions: conv,
			Spend:       spend,
			Revenue:     rev,
		})
	}
	return campaigns, nil
}

// ── Aggregates ─────────────────────────────────────────────────────────────

func totalMetrics(all []Metrics) Metrics {
	var t Campaign
	t.Name = "TOTAL"
	for _, m := range all {
		t.Impressions += m.Impressions
		t.Clicks += m.Clicks
		t.Conversions += m.Conversions
		t.Spend += m.Spend
		t.Revenue += m.Revenue
	}
	return computeMetrics(t)
}

// ── Display ────────────────────────────────────────────────────────────────

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	gray   = "\033[90m"
)

func colorROAS(v float64) string {
	s := fmt.Sprintf("%.2fx", v)
	switch {
	case v >= 4:
		return green + bold + s + reset
	case v >= 2:
		return yellow + s + reset
	default:
		return red + s + reset
	}
}

func printTable(metrics []Metrics) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "%s%-24s\tImpr\tClicks\tConv\tCTR%%\tCVR%%\tCPC\tCPA\tROAS\tSpend\tRevenue%s\n",
		bold+cyan, "Campaign", reset)
	fmt.Fprintf(w, "%s%-24s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s%s\n",
		gray, strings.Repeat("─", 24),
		"────────", "────────", "────────",
		"────────", "────────", "────────",
		"────────", "────────", "────────", "────────", reset)

	for _, m := range metrics {
		fmt.Fprintf(w, "%-24s\t%.0f\t%.0f\t%.0f\t%.2f%%\t%.2f%%\t$%.2f\t$%.2f\t%s\t$%.2f\t$%.2f\n",
			m.Name, m.Impressions, m.Clicks, m.Conversions,
			m.CTR, m.CVR, m.CPC, m.CPA,
			colorROAS(m.ROAS),
			m.Spend, m.Revenue)
	}

	w.Flush()
}

func printSummary(total Metrics) {
	fmt.Printf("\n%s%s── Summary ──────────────────────────────────────%s\n", bold, cyan, reset)
	fmt.Printf("  %-22s %s$%.2f%s\n", "Total Spend:", bold, total.Spend, reset)
	fmt.Printf("  %-22s %s$%.2f%s\n", "Total Revenue:", bold, total.Revenue, reset)
	fmt.Printf("  %-22s %s%.2fx%s\n", "Blended ROAS:", bold, total.ROAS, reset)
	fmt.Printf("  %-22s %s%.2f%%%s\n", "Blended CTR:", bold, total.CTR, reset)
	fmt.Printf("  %-22s %s%.2f%%%s\n", "Blended CVR:", bold, total.CVR, reset)
	fmt.Printf("  %-22s %s$%.2f%s\n", "Blended CPA:", bold, total.CPA, reset)
	fmt.Printf("  %-22s %s$%.2f%s\n", "Blended CPC:", bold, total.CPC, reset)
	fmt.Printf("%s%s─────────────────────────────────────────────────%s\n\n", bold, cyan, reset)
}

// ── Main ───────────────────────────────────────────────────────────────────

func main() {
	filePath := flag.String("file", "", "Path to campaign CSV file (required)")
	sortBy := flag.String("sort", "roas", "Sort by: roas | ctr | cvr | cpa | spend | revenue")
	minROAS := flag.Float64("min-roas", 0, "Filter: show only campaigns with ROAS >= value")
	flag.Parse()

	if *filePath == "" {
		fmt.Fprintln(os.Stderr, red+"Error: --file flag is required."+reset)
		fmt.Fprintln(os.Stderr, gray+"Usage: go run main.go --file data/sample_campaigns.csv"+reset)
		os.Exit(1)
	}

	campaigns, err := loadCSV(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, red+"Error: %v\n"+reset, err)
		os.Exit(1)
	}

	if len(campaigns) == 0 {
		fmt.Fprintln(os.Stderr, yellow+"Warning: no campaign rows found in file."+reset)
		os.Exit(0)
	}

	// Compute metrics
	all := make([]Metrics, 0, len(campaigns))
	for _, c := range campaigns {
		m := computeMetrics(c)
		if m.ROAS >= *minROAS {
			all = append(all, m)
		}
	}

	// Sort
	sort.Slice(all, func(i, j int) bool {
		switch strings.ToLower(*sortBy) {
		case "ctr":
			return all[i].CTR > all[j].CTR
		case "cvr":
			return all[i].CVR > all[j].CVR
		case "cpa":
			return all[i].CPA < all[j].CPA // lower is better
		case "spend":
			return all[i].Spend > all[j].Spend
		case "revenue":
			return all[i].Revenue > all[j].Revenue
		default:
			return all[i].ROAS > all[j].ROAS
		}
	})

	total := totalMetrics(all)

	// Print header
	fmt.Printf("\n%s%s Marketing Metrics CLI%s — %s%s%s\n",
		bold, cyan, reset, gray, *filePath, reset)
	fmt.Printf("%s%d campaigns · sorted by %s%s\n\n",
		gray, len(all), strings.ToLower(*sortBy), reset)

	printTable(all)
	printSummary(total)

	// Best / worst
	if len(all) > 1 {
		fmt.Printf("%sBest ROAS  :%s %s (%.2fx)\n", green, reset, all[0].Name, all[0].ROAS)
		fmt.Printf("%sLowest CPA :%s %s ($%.2f)\n\n", yellow, reset,
			func() (string, float64) {
				best := all[0]
				for _, m := range all[1:] {
					if m.Conversions > 0 && m.CPA < best.CPA {
						best = m
					}
				}
				return best.Name, best.CPA
			}())
	}
}
