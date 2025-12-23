package date

import (
	"encoding/json"
	"time"
)

// Date represents a calendar date.
//
// It is primarily used for JSON request/response payloads where callers may send
// a date-only string (YYYY-MM-DD). Internally it is stored as a time.Time at
// midnight UTC.
//
// It also accepts RFC3339 timestamps for flexibility.
type Date struct {
	time.Time
}

// UnmarshalJSON accepts either:
// - null
// - an empty string (treated as zero)
// - a date-only string in YYYY-MM-DD
// - a RFC3339 timestamp
func (d *Date) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil {
		d.Time = time.Time{}
		return nil
	}
	if *s == "" {
		d.Time = time.Time{}
		return nil
	}

	// Date-only format.
	if len(*s) == len("2006-01-02") {
		t, err := time.ParseInLocation("2006-01-02", *s, time.UTC)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	}

	// Fallback to timestamp.
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}
