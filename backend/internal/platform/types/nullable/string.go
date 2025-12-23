// Package nullable provides types that properly marshal NULL database values to JSON.
package nullable

import (
	"database/sql"
	"encoding/json"
)

// String is a nullable string that marshals to JSON as a string or null.
// It wraps sql.NullString to provide proper JSON serialization.
type String struct {
	sql.NullString
}

// NewString creates a new nullable string.
func NewString(s string) String {
	return String{
		NullString: sql.NullString{
			String: s,
			Valid:  s != "",
		},
	}
}

// NewStringFromPtr creates a new nullable string from a pointer.
func NewStringFromPtr(s *string) String {
	if s == nil {
		return String{
			NullString: sql.NullString{
				Valid: false,
			},
		}
	}
	return NewString(*s)
}

// MarshalJSON implements json.Marshaler.
// It marshals the string value or null if invalid.
func (ns String) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON implements json.Unmarshaler.
// It unmarshals a JSON string or null to the nullable string.
func (ns *String) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil {
		ns.Valid = false
		ns.String = ""
		return nil
	}
	ns.Valid = true
	ns.String = *s
	return nil
}

// Ptr returns a pointer to the string value, or nil if invalid.
func (ns String) Ptr() *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

// ValueOrDefault returns the string value or a default if invalid.
func (ns String) ValueOrDefault(defaultValue string) string {
	if !ns.Valid {
		return defaultValue
	}
	return ns.String
}
