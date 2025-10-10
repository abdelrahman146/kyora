package handlers

import (
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type ForgotPasswordHandler struct {
}

func NewForgotPasswordHandler() *ForgotPasswordHandler {
	return &ForgotPasswordHandler{}
}

func (h *ForgotPasswordHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/forgot-password", h.Index)
}

func (h *ForgotPasswordHandler) Index(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Forgot Password",
		Description: "Forgot Password page",
		Keywords:    "forgot password, Kyora",
		Path:        "/forgot-password",
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ForgotPassword())
}
