package webcontext

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/types"
	"github.com/spf13/viper"
)

var (
	PageInfoKey = types.ContextKey{Name: "page_info"}
)

type PageInfo struct {
	Locale      string
	Dir         string
	Title       string
	Description string
	Keywords    string
	Path        string
	Breadcrumbs []Breadcrumb
}

type Breadcrumb struct {
	Href  string
	Label string
}

func SetupPageInfo(ctx context.Context, info PageInfo) context.Context {
	return context.WithValue(ctx, PageInfoKey, info)
}

func GetPageTitle(ctx context.Context) string {
	info := GetVal[PageInfo](ctx, PageInfoKey)
	return info.Title
}

func GetPageDescription(ctx context.Context) string {
	info := GetVal[PageInfo](ctx, PageInfoKey)
	return info.Description
}

func GetPageKeywords(ctx context.Context) string {
	info := GetVal[PageInfo](ctx, PageInfoKey)
	return info.Keywords
}

func GetPageBreadcrumbs(ctx context.Context) []Breadcrumb {
	info := GetVal[PageInfo](ctx, PageInfoKey)
	return info.Breadcrumbs
}

func GetActivePath(ctx context.Context) string {
	info := GetVal[PageInfo](ctx, PageInfoKey)
	return info.Path
}

func GetPageLocale(ctx context.Context) string {
	info := GetVal[PageInfo](ctx, PageInfoKey)
	if info.Locale == "" {
		return "en"
	}
	return info.Locale
}

func GetPageDir(ctx context.Context) string {
	info := GetVal[PageInfo](ctx, PageInfoKey)
	if info.Dir == "" {
		return "ltr"
	}
	return info.Dir
}

func ComposeFullTitle(pageTitle string) string {
	siteName := viper.GetString("site.name")
	if pageTitle == "" {
		return siteName
	}
	if siteName == "" {
		return pageTitle
	}
	return pageTitle + " - " + siteName
}
