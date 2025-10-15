package webutils

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/utils"
)

func NewFromDate(ctx context.Context, dateStr string) time.Time {
	from := time.Now().AddDate(0, 0, -30)
	if dateStr != "" {
		parsed, err := time.Parse(time.RFC3339, dateStr)
		if err == nil {
			from = parsed
		} else {
			utils.Log.FromContext(ctx).Warn("failed to parse 'from' date", "error", err, "dateStr", dateStr)
		}
	}
	return from
}

func NewToDate(ctx context.Context, dateStr string) time.Time {
	to := time.Now()
	if dateStr != "" {
		parsed, err := time.Parse(time.RFC3339, dateStr)
		if err == nil {
			to = parsed
		} else {
			utils.Log.FromContext(ctx).Warn("failed to parse 'to' date", "error", err, "dateStr", dateStr)
		}
	}
	return to
}

var DateOptions = []string{
	"Last 7 days",
	"Last 30 days",
	"Last 90 days",
	"This month",
	"Last month",
	"This year",
	"Last year",
	"All time",
	"Custom range",
}

func GetDateRange(option string) (time.Time, time.Time) {
	now := time.Now()
	var from, to time.Time

	switch option {
	case "Last 7 days":
		from = now.AddDate(0, 0, -7)
		to = now
	case "Last 30 days":
		from = now.AddDate(0, 0, -30)
		to = now
	case "Last 90 days":
		from = now.AddDate(0, 0, -90)
		to = now
	case "This month":
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		to = now
	case "Last month":
		firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		from = firstOfThisMonth.AddDate(0, -1, 0)
		to = firstOfThisMonth.AddDate(0, 0, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
	case "This year":
		from = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		to = now
	case "Last year":
		from = time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
		to = time.Date(now.Year()-1, 12, 31, 23, 59, 59, 0, now.Location())
	case "All time":
		from = time.Time{} // Zero time
		to = now
	default:
		from = now.AddDate(0, 0, -30)
		to = now
	}
	return from, to
}
