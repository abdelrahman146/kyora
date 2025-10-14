package utils

import sd "github.com/shopspring/decimal"

func init() {
	// Emit decimals as JSON numbers instead of quoted strings
	sd.MarshalJSONWithoutQuotes = true
}
