package timeseries

import (
	"fmt"
	"time"
)

type Granularity int

const (
	Hourly Granularity = iota
	Daily
	Weekly
	Monthly
	Quarterly
	Yearly
)

func (g Granularity) String() string {
	switch g {
	case Hourly:
		return "hourly"
	case Daily:
		return "daily"
	case Weekly:
		return "weekly"
	case Monthly:
		return "monthly"
	case Quarterly:
		return "quarterly"
	case Yearly:
		return "yearly"
	default:
		return "unknown"
	}
}

func (g Granularity) Bucket() string {
	switch g {
	case Hourly:
		return "hour"
	case Daily:
		return "day"
	case Weekly:
		return "week"
	case Monthly:
		return "month"
	case Quarterly:
		return "quarter"
	case Yearly:
		return "year"
	default:
		return "day"
	}
}

type TimeSeries struct {
	Granularity Granularity     `json:"granularity"`
	Series      []TimeSeriesRow `json:"series"`
}

// TimeSeriesRow is a lightweight row for time series aggregations.
type TimeSeriesRow struct {
	Timestamp time.Time `gorm:"column:timestamp" json:"timestamp"`
	Label     string    `gorm:"-" json:"label"`
	Value     float64   `gorm:"column:value" json:"value"`
}

// New creates a time series with a sensible granularity based on the span
// between 'from' and 'to'. It aligns the start to a natural boundary and fills in
// human-friendly tick labels and corresponding values. Missing data points are
// filled with zero values. It assumes Monday as the start of the week.
//
// Examples of outputs by span:
//   - <= 48h: "15:00 2 Jan", "16:00 2 Jan", ... (Hourly)
//   - <= 45d: "Mon 2 Jan", "Tue 3 Jan", ... (Daily)
//   - <= 120d: "Wk2 Jan 2025", ... (Weekly; week-of-month)
//   - <= 2y:  "Jan 2025", "Feb 2025", ... (Monthly)
//   - <= 5y:  "Q1 2025", "Q2 2025", ... (Quarterly)
//   - > 5y:   "2025", "2026", ... (Yearly)
func New(rows []TimeSeriesRow, granularity Granularity) *TimeSeries {
	timeSeries := &TimeSeries{
		Granularity: granularity,
	}
	if len(rows) == 0 {
		return nil
	}
	includeYear := rows[0].Timestamp.Year() != rows[len(rows)-1].Timestamp.Year()
	for _, r := range rows {
		timeSeries.Series = append(timeSeries.Series, TimeSeriesRow{
			Timestamp: r.Timestamp,
			Label:     formatLabel(r.Timestamp, granularity, includeYear),
			Value:     r.Value,
		})
	}
	return timeSeries
}

// --- helpers ---

func GetTimeGranularityByDateRange(from, to time.Time) Granularity {
	d := to.Sub(from)
	switch {
	case d <= 48*time.Hour:
		return Hourly
	case d <= 45*24*time.Hour:
		return Daily
	case d <= 120*24*time.Hour:
		return Weekly
	case d <= 24*30*24*time.Hour: // ~2 years
		return Monthly
	case d <= 60*30*24*time.Hour: // ~5 years
		return Quarterly
	default:
		return Yearly
	}
}

func formatLabel(t time.Time, g Granularity, includeYear bool) string {
	switch g {
	case Hourly:
		// Always include day+month to avoid ambiguity across days
		return t.Format("15:04 2 Jan")
	case Daily:
		// Show weekday when zoomed into days for readability
		// If the overall span crosses years, include the year
		if includeYear {
			return t.Format("Mon 2 Jan 2006")
		}
		return t.Format("Mon 2 Jan")
	case Weekly:
		wk := weekOfMonth(t, time.Monday)
		if includeYear {
			return fmt.Sprintf("Wk%d %s", wk, t.Format("Jan 2006"))
		}
		return fmt.Sprintf("Wk%d %s", wk, t.Format("Jan"))
	case Monthly:
		if includeYear {
			return t.Format("Jan 2006")
		}
		return t.Format("Jan")
	case Quarterly:
		_, q := quarterOf(t)
		return fmt.Sprintf("Q%d %d", q, t.Year())
	case Yearly:
		return t.Format("2006")
	default:
		return t.String()
	}
}

func weekOfMonth(t time.Time, weekStart time.Weekday) int {
	first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	shift := (int(first.Weekday()) - int(weekStart) + 7) % 7
	return ((t.Day() - 1 + shift) / 7) + 1
}

func quarterOf(t time.Time) (year int, q int) {
	m := int(t.Month())
	q = (m-1)/3 + 1
	return t.Year(), q
}
