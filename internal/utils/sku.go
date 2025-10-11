package utils

import (
	"regexp"
	"strings"
)

// skuHelper provides helpers to generate human-friendly SKU codes.
// Format: XXX-XXX-XXXXXX where
// - first segment: 3-char product code from product name
// - second segment: 3-char variant code from variant name
// - third segment: random uppercase alphanumeric (ULID fragment)
// Optionally, callers can prepend a store-specific code if desired.
type skuHelper struct{}

var nonAlphaNum = regexp.MustCompile(`[^A-Za-z0-9]+`)

// Generate builds a SKU from product and variant names and a random suffix.
// storeID is accepted to make it easier to evolve the format later, but is not encoded directly.
func (skuHelper) Generate(storeID, productName, variantName string) string {
	p := codeFromName(productName)
	v := codeFromName(variantName)
	// Use ULID fragment as random uppercase alphanumeric suffix
	rand := ID.NewBase62(6)
	return strings.ToUpper(p + "-" + v + "-" + rand)
}

func codeFromName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "XXX"
	}
	// Remove non-alphanumeric and split by spaces after cleanup
	cleaned := nonAlphaNum.ReplaceAllString(s, " ")
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		return "XXX"
	}
	parts := strings.Fields(cleaned)
	// Prefer first token; fall back to entire cleaned string
	token := parts[0]
	token = nonAlphaNum.ReplaceAllString(token, "")
	if len(token) >= 3 {
		return strings.ToUpper(token[:3])
	}
	// Pad to 3 chars using X if shorter
	token = strings.ToUpper(token)
	for len(token) < 3 {
		token += "X"
	}
	return token
}

var SKU = skuHelper{}
