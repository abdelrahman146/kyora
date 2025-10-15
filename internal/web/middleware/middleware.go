package middleware

import "github.com/abdelrahman146/kyora/internal/types"

var (
	ClaimsKey     = "claims"
	UserKey       = "user"
	loginPath     = "/login"
	StoreKey      = types.ContextKey{Name: "store"}
	StoresListKey = types.ContextKey{Name: "stores_list"}
)
