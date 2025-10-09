package webutils

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, status int, content templ.Component) {
	w := c.Writer
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Status(status)
	if err := content.Render(c.Request.Context(), c.Writer); err != nil {
		utils.Log.FromContext(c.Request.Context()).Error("failed to render content", "error", err)
	}
}

func RenderFragments(c *gin.Context, status int, content templ.Component, keys ...webcontext.FragmentKey) {
	w := c.Writer
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Status(status)
	ids := make([]any, len(keys))
	for i, key := range keys {
		ids[i] = key
	}
	if err := templ.RenderFragments(c.Request.Context(), w, content, ids...); err != nil {
		utils.Log.FromContext(c.Request.Context()).Error("failed to render fragments", "error", err)
	}
}

func Redirect(c *gin.Context, location string) {
	c.Header("HX-Redirect", location)
	c.Status(http.StatusOK)
}
