package handlers

import (
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type ResetPasswordHandler struct {
}

func NewResetPasswordHandler() *ResetPasswordHandler {
	return &ResetPasswordHandler{}
}

func (h *ResetPasswordHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/reset-password", h.Index)
}

func (h *ResetPasswordHandler) Index(c *gin.Context) {
	info := webcontext.PageInfo{
		Locale:      "en",
		Dir:         "ltr",
		Title:       "Reset Password",
		Description: "Reset Password page",
		Keywords:    "reset password, Kyora",
		Path:        "/reset-password",
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, pages.ResetPassword())
}
