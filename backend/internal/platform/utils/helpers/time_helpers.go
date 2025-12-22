package helpers

import (
	"math"
	"time"
)

// CeilPositiveHoursUntil returns the ceiling of the remaining hours until t.
// It clamps negative/zero durations to 0.
func CeilPositiveHoursUntil(t time.Time) int {
	d := time.Until(t)
	if d <= 0 {
		return 0
	}
	return int(math.Ceil(d.Hours()))
}

// CeilPositiveDaysUntil returns the ceiling of the remaining days until t.
// It clamps negative/zero durations to 0.
func CeilPositiveDaysUntil(t time.Time) int {
	d := time.Until(t)
	if d <= 0 {
		return 0
	}
	return int(math.Ceil(d.Hours() / 24))
}
