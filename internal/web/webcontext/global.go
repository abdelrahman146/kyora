package webcontext

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/spf13/viper"
)

type ContextKey struct {
	Name string
}

func GetAppName(ctx context.Context) string {
	return GetValFromConfig[string](ctx, "app.name")
}

func GetBaseURL(ctx context.Context) string {
	return GetValFromConfig[string](ctx, "site.base_url")
}

func GetSiteName(ctx context.Context) string {
	return GetValFromConfig[string](ctx, "site.name")
}

func GetSiteDescription(ctx context.Context) string {
	return GetValFromConfig[string](ctx, "site.description")
}

func GetValFromConfig[T any](ctx context.Context, key string) T {
	var val T
	if err := viper.UnmarshalKey(key, &val); err != nil {
		utils.Log.FromContext(ctx).Warn("failed to get value from config", "key", key, "error", err)
		return *new(T)
	}
	return val
}

func GetVal[T any](ctx context.Context, key ContextKey) T {
	val, ok := ctx.Value(key).(T)
	if !ok {
		utils.Log.FromContext(ctx).Warn("failed to get value from context", "key", key.Name)
		return *new(T)
	}
	return val
}
