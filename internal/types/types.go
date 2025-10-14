package types

import (
	"time"
)

// TimeSeriesRow is a lightweight row for time series aggregations.
type TimeSeriesRow struct {
	Timestamp time.Time `gorm:"column:timestamp"`
	Value     float64   `gorm:"column:value"`
}

// KeyValue is a simple key-value row for group-by aggregations.
type KeyValue struct {
	Key   string  `gorm:"column:key"`
	Value float64 `gorm:"column:value"`
}
