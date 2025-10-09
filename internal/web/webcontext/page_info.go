package webcontext

import (
	"context"

	"github.com/spf13/viper"
)

var (
	PageInfoKey = ContextKey{"page_info"}
)

type PageInfo struct {
	Title       string
	Description string
	Keywords    string
	Breadcrumbs []Breadcrumb
}

type Breadcrumb struct {
	Title string
	Link  string
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
