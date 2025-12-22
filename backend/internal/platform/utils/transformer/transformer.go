package transformer

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

func ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func FromNullString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

func ToNullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func FromNullInt64(ni sql.NullInt64) int64 {
	if !ni.Valid {
		return 0
	}
	return ni.Int64
}

func ToNullFloat64(f float64) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}

func FromNullFloat64(nf sql.NullFloat64) float64 {
	if !nf.Valid {
		return 0
	}
	return nf.Float64
}

func ToNullBool(b bool) sql.NullBool {
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}

func FromNullBool(nb sql.NullBool) bool {
	if !nb.Valid {
		return false
	}
	return nb.Bool
}

func ToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func FromNullTime(nt sql.NullTime) sql.NullTime {
	if !nt.Valid {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  nt.Time,
		Valid: true,
	}
}

func ToNullDecimal(d decimal.Decimal) sql.NullString {
	if d.IsZero() {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: d.String(),
		Valid:  true,
	}
}

func FromNullDecimal(nd decimal.NullDecimal) decimal.Decimal {
	if !nd.Valid {
		return decimal.Zero
	}
	return nd.Decimal
}
