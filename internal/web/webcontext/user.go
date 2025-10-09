package webcontext

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/domain/account"
)

var (
	UserKey = ContextKey{"user"}
)

func SetupUserContext(ctx context.Context, user *account.User) context.Context {
	ctx = context.WithValue(ctx, UserKey, user)
	return ctx
}

func IsAuthenticated(ctx context.Context) bool {
	return GetUserFromContext(ctx) != nil
}

func GetUserFromContext(ctx context.Context) *account.User {
	return GetVal[*account.User](ctx, UserKey)
}
