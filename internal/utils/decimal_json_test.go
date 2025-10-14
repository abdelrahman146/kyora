package utils

import (
	"encoding/json"
	"testing"

	sd "github.com/shopspring/decimal"
)

func TestDecimalJSONNumber(t *testing.T) {
	// init() in decimal_json.go sets MarshalJSONWithoutQuotes = true
	d, err := sd.NewFromString("12.34")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(b) != "12.34" {
		t.Fatalf("expected numeric JSON without quotes, got %s", string(b))
	}
	var d2 sd.Decimal
	if err := json.Unmarshal(b, &d2); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !d.Equal(d2) {
		t.Fatalf("round-trip mismatch: %s vs %s", d.String(), d2.String())
	}
}

func TestDecimalScanValue(t *testing.T) {
	var d sd.Decimal
	if err := d.Scan("45.67"); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if d.String() != "45.67" {
		t.Fatalf("unexpected: %s", d.String())
	}
	val, err := d.Value()
	if err != nil {
		t.Fatalf("value: %v", err)
	}
	if s, ok := val.(string); !ok || s != "45.67" {
		t.Fatalf("unexpected driver.Value: %#v", val)
	}
}
