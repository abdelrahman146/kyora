package analytics

import (
	"context"
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

type TimeSeries struct {
	Granularity Granularity
	Labels      []string
	Values      []float64
}

type TimeSeriesRow struct {
	Timestamp time.Time
	Value     float64
}

// newTimeSeries creates a time series with a sensible granularity based on the span
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
func newTimeSeries(ctx context.Context, rows []TimeSeriesRow, from, to time.Time) *TimeSeries {
	if len(rows) == 0 {
		return &TimeSeries{
			Granularity: Daily,
			Labels:      []string{},
			Values:      []float64{},
		}
	}

	g := chooseGranularity(from, to)
	labels := make([]string, 0, 64)
	values := make([]float64, 0, 64)

	loc := from.Location()
	from = from.In(loc)
	to = to.In(loc)

	start := alignStart(from, g)
	end := alignStart(to, g)

	rowIdx := 0
	for t := start; !t.After(end); t = nextTick(t, g) {
		labels = append(labels, formatLabel(t, g, from.Year() != to.Year()))
		if rowIdx < len(rows) && rows[rowIdx].Timestamp.Equal(t) {
			values = append(values, rows[rowIdx].Value)
			rowIdx++
		} else {
			values = append(values, 0)
		}
	}

	return &TimeSeries{
		Granularity: g,
		Labels:      labels,
		Values:      values,
	}
}

// generateTimeSeriesLabels chooses a sensible granularity based on the span between
// from and to, snaps the start to a natural boundary, and returns human-friendly
// tick labels for charts. It assumes Monday as the start of the week.
//
// Examples of outputs by span:
//   - <= 48h: "15:00 2 Jan", "16:00 2 Jan", ... (Hourly)
//   - <= 45d: "Mon 2 Jan", "Tue 3 Jan", ... (Daily)
//   - <= 120d: "Wk2 Jan 2025", ... (Weekly; week-of-month)
//   - <= 2y:  "Jan 2025", "Feb 2025", ... (Monthly)
//   - <= 5y:  "Q1 2025", "Q2 2025", ... (Quarterly)
//   - > 5y:   "2025", "2026", ... (Yearly)
func (s *analyticsService) generateTimeSeriesLabels(from, to time.Time) []string {
	if to.Before(from) {
		from, to = to, from
	}

	// Normalize nanoseconds and use the same location throughout
	loc := from.Location()
	from = from.In(loc)
	to = to.In(loc)

	g := chooseGranularity(from, to)
	start := alignStart(from, g)
	labels := make([]string, 0, 64)

	includeYear := from.Year() != to.Year()

	for t := start; !t.After(to); t = nextTick(t, g) {
		labels = append(labels, formatLabel(t, g, includeYear))
	}
	return dedupeLeading(labels)
}

// --- helpers ---

func chooseGranularity(from, to time.Time) Granularity {
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

func alignStart(t time.Time, g Granularity) time.Time {
	switch g {
	case Hourly:
		return t.Truncate(time.Hour)
	case Daily:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case Weekly:
		return startOfWeek(t, time.Monday)
	case Monthly:
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case Quarterly:
		y, q := quarterOf(t)
		return startOfQuarter(y, q, t.Location())
	case Yearly:
		return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	default:
		return t
	}
}

func nextTick(t time.Time, g Granularity) time.Time {
	switch g {
	case Hourly:
		return t.Add(time.Hour)
	case Daily:
		return t.AddDate(0, 0, 1)
	case Weekly:
		return t.AddDate(0, 0, 7)
	case Monthly:
		return t.AddDate(0, 1, 0)
	case Quarterly:
		return t.AddDate(0, 3, 0)
	case Yearly:
		return t.AddDate(1, 0, 0)
	default:
		return t
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

func startOfWeek(t time.Time, weekStart time.Weekday) time.Time {
	wd := int(t.Weekday())
	ws := int(weekStart)
	// Go's Weekday: Sunday=0 ... Monday=1 ... Saturday=6
	// Convert to offset back to weekStart
	delta := (7 + wd - ws) % 7
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return midnight.AddDate(0, 0, -delta)
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

func startOfQuarter(year, q int, loc *time.Location) time.Time {
	month := time.Month((q-1)*3 + 1)
	return time.Date(year, month, 1, 0, 0, 0, 0, loc)
}

// dedupeLeading removes any duplicate leading labels that may occur when the
// aligned start is before the original 'from' but produces the same first label
// as the next tick. Keeps ordering stable.
func dedupeLeading(labels []string) []string {
	if len(labels) < 2 {
		return labels
	}
	out := make([]string, 0, len(labels))
	prev := ""
	for i, l := range labels {
		if i == 0 || l != prev {
			out = append(out, l)
			prev = l
		}
	}
	return out
}
