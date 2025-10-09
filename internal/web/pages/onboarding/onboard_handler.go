package onboardingPage

import (
	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/web/middleware"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"

	"github.com/gin-gonic/gin"
)

type OnboardHandler struct {
	onboardService *account.OnboardingService
}

func NewOnboardHandler(onboardService *account.OnboardingService) *OnboardHandler {
	return &OnboardHandler{onboardService}
}

func (h *OnboardHandler) RegisterRoutes(r gin.IRoutes) {
	r.GET("/onboard", middleware.GuestRequiredMiddleware, h.Index)
}

func (h *OnboardHandler) Index(c *gin.Context) {
	info := webcontext.PageInfo{
		Title:       "Onboarding",
		Description: "Onboarding page ",
		Keywords:    "onboarding, zard",
	}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, 200, Onboarding())
}
