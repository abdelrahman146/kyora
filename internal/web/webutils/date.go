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
