# Marketing Metrics CLI — Go / CLI de Metricas de Marketing en Go

A command-line tool that reads campaign performance data from a CSV file and computes standard digital marketing KPIs — built entirely with Go's standard library, no external dependencies.

Una herramienta de linea de comandos que lee datos de campanas de marketing desde un CSV y calcula KPIs de marketing digital — construida con la libreria estandar de Go, sin dependencias externas.

---

## Key Metrics / Metricas Calculadas

| KPI | Formula | Description |
|-----|---------|-------------|
| **CTR** | clicks / impressions × 100 | Click-Through Rate |
| **CVR** | conversions / clicks × 100 | Conversion Rate |
| **CPC** | spend / clicks | Cost Per Click |
| **CPA** | spend / conversions | Cost Per Acquisition |
| **ROAS** | revenue / spend | Return on Ad Spend |
| **Rev/Conv** | revenue / conversions | Revenue Per Conversion |

---

## Usage / Uso

```bash
# Run with sample data
go run main.go --file data/sample_campaigns.csv

# Sort by CTR instead of ROAS (default)
go run main.go --file data/sample_campaigns.csv --sort ctr

# Filter: only show campaigns with ROAS >= 3
go run main.go --file data/sample_campaigns.csv --min-roas 3.0

# Build binary and run
go build -o mkt-cli
./mkt-cli --file data/sample_campaigns.csv --sort revenue
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--file` | required | Path to CSV file |
| `--sort` | `roas` | Sort by: `roas`, `ctr`, `cvr`, `cpa`, `spend`, `revenue` |
| `--min-roas` | `0` | Filter campaigns below this ROAS threshold |

---

## CSV Format / Formato del CSV

```csv
campaign_name,impressions,clicks,conversions,spend,revenue
Google Search - Brand,120000,4800,320,1850.00,9200.00
Meta - Retargeting,85000,2100,210,1200.00,6300.00
```

Required columns: `campaign_name`, `impressions`, `clicks`, `conversions`, `spend`, `revenue`

---

## Sample Output / Ejemplo de Salida

```
 Marketing Metrics CLI — data/sample_campaigns.csv
 10 campaigns · sorted by roas

 Campaign                   Impr      Clicks   Conv   CTR%    CVR%    CPC     CPA     ROAS    Spend      Revenue
 ────────────────────────   ────────  ──────   ────   ──────  ──────  ──────  ──────  ──────  ─────────  ─────────
 Email - Newsletter              0   12000     680    0.00%   5.67%   $0.03   $0.47   25.50x  $320.00    $8160.00
 Google Search - Brand      120000    4800     320    4.00%   6.67%   $0.39   $5.78    4.97x  $1850.00   $9200.00
 Meta - Retargeting          85000    2100     210    2.47%  10.00%   $0.57   $5.71    5.25x  $1200.00   $6300.00
 ...

── Summary ──────────────────────────────────────
  Total Spend:           $16270.00
  Total Revenue:         $44835.00
  Blended ROAS:          2.76x
  Blended CTR:           1.92%
  Blended CVR:           4.67%
  Blended CPA:           $9.00
─────────────────────────────────────────────────

Best ROAS  : Email - Newsletter (25.50x)
Lowest CPA : Email - Newsletter ($0.47)
```

---

## Why Go / Por que Go

Python is great for data analysis, but Go offers real advantages for CLI tooling:

- **Performance:** Processes large CSVs (millions of rows) significantly faster than pandas
- **Concurrency:** Goroutines enable parallel processing of multiple files natively  
- **Single binary:** Compiles to a standalone executable — no runtime, no venv, no pip
- **Type safety:** Compile-time checks catch errors before runtime

---

## Project Structure / Estructura del Proyecto

```
Marketing-Metrics-CLI-Go/
├── main.go                   # CLI logic + KPI computation
├── go.mod                    # Module definition
└── data/
    └── sample_campaigns.csv  # Sample dataset (10 campaigns)
```

---

## Tech Stack

- **Language:** Go 1.21
- **Dependencies:** None (standard library only)
- **Packages used:** `encoding/csv`, `flag`, `text/tabwriter`, `sort`

---

## Author / Autor

**Julian David Urrego Lancheros**  
Data Analyst · Python Developer · Marketing Science  
[GitHub](https://github.com/JulianDataScienceExplorerV2) · [juliandavidurrego2011@gmail.com](mailto:juliandavidurrego2011@gmail.com)
