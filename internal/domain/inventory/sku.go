package inventory

import (
	"regexp"
	"strings"

	"github.com/abdelrahman146/kyora/internal/utils"
)

var nonAlphaNum = regexp.MustCompile(`[^A-Za-z0-9]+`)

// GenerateSku generates human-friendly SKU codes.
// New standardized format: STORE-PRO-VAR-RAND where
// - STORE: up to 6-char sanitized store code (alphanumeric, uppercased)
// - PRO: 3-char product code from product name
// - VAR: 3-char variant code from variant name
// - RAND: random 6-char uppercase alphanumeric (Base62 fragment)
func GenerateSku(storeCode, productName, variantName string) string {
	s := sanitizeStoreCode(storeCode)
	p := codeFromName(productName)
	v := codeFromName(variantName)
	// Use ULID fragment as random uppercase alphanumeric suffix
	rand := utils.ID.NewBase62(6)
	return strings.ToUpper(s + "-" + p + "-" + v + "-" + rand)
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

// sanitizeStoreCode normalizes the store code for inclusion in the SKU.
// - trims spaces
// - strips non-alphanumeric characters
// - uppercases
// - limits length to 6 characters
// - falls back to "STR" if empty
func sanitizeStoreCode(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "STR"
	}
	cleaned := nonAlphaNum.ReplaceAllString(s, "")
	if cleaned == "" {
		return "STR"
	}
	if len(cleaned) > 6 {
		cleaned = cleaned[:6]
	}
	return strings.ToUpper(cleaned)
}
