package utils

import (
	"database/sql"
	"time"
)

type transformHelper struct{}

func (transformHelper) ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func (transformHelper) ToPtrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (transformHelper) ToPtrBool(b bool) *bool {
	return &b
}

func (transformHelper) ToPtrInt(i int) *int {
	return &i
}

func (transformHelper) ToPtrFloat64(f float64) *float64 {
	return &f
}

func (transformHelper) ToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func (transformHelper) ToTime(nt sql.NullTime) time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return nt.Time
}

func (transformHelper) ToString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

func (transformHelper) ToBool(nb sql.NullBool) bool {
	if !nb.Valid {
		return false
	}
	return nb.Bool
}

func (transformHelper) ToInt(ni sql.NullInt64) int {
	if !ni.Valid {
		return 0
	}
	return int(ni.Int64)
}

func (transformHelper) ToFloat64(nf sql.NullFloat64) float64 {
	if !nf.Valid {
		return 0
	}
	return nf.Float64
}

var Transform = &transformHelper{}
