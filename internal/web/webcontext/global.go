package webcontext

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/spf13/viper"
)

type ContextKey struct {
	Name string
}

var (
	AppNameKey         = ContextKey{"appName"}
	BaseURLKey         = ContextKey{"baseURL"}
	SiteNameKey        = ContextKey{"siteName"}
	SiteDescriptionKey = ContextKey{"siteDescription"}
)

func SetupGlobalContext(ctx context.Context) context.Context {
	appName := viper.GetString("app.name")
	baseURL := viper.GetString("site.base_url")
	siteName := viper.GetString("site.title")
	siteDescription := viper.GetString("site.description")
	ctx = context.WithValue(ctx, AppNameKey, appName)
	ctx = context.WithValue(ctx, BaseURLKey, baseURL)
	ctx = context.WithValue(ctx, SiteNameKey, siteName)
	ctx = context.WithValue(ctx, SiteDescriptionKey, siteDescription)
	return ctx
}

func GetVal[T any](ctx context.Context, key ContextKey) T {
	val, ok := ctx.Value(key).(T)
	if !ok {
		utils.Log.FromContext(ctx).Warn("failed to get value from context", "key", key.Name)
		return *new(T)
	}
	return val
}

func GetAppName(ctx context.Context) string {
	return GetVal[string](ctx, AppNameKey)
}

func GetBaseURL(ctx context.Context) string {
	return GetVal[string](ctx, BaseURLKey)
}

func GetSiteName(ctx context.Context) string {
	return GetVal[string](ctx, SiteNameKey)
}

func GetSiteDescription(ctx context.Context) string {
	return GetVal[string](ctx, SiteDescriptionKey)
}
