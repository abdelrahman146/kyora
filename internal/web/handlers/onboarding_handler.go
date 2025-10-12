package handlers

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/abdelrahman146/kyora/internal/web/views/pages"
	"github.com/abdelrahman146/kyora/internal/web/webcontext"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
)

type OnboardingHandler struct {
	onboarding *account.OnboardingService
	auth       *account.AuthenticationService
}

const (
	onboardingPath = "/onboarding"
)

func NewOnboardingHandler(onboarding *account.OnboardingService, auth *account.AuthenticationService) *OnboardingHandler {
	return &OnboardingHandler{onboarding: onboarding, auth: auth}
}

func (h *OnboardingHandler) RegisterRoutes(r gin.IRoutes) {
	// wizard
	r.GET(onboardingPath, h.Index)
	r.POST(onboardingPath+"/step2", h.Step2)
	r.POST(onboardingPath+"/complete", h.Complete)
	r.GET(onboardingPath+"/slug-availability", h.SlugAvailability)
}

func (h *OnboardingHandler) Index(c *gin.Context) {
	info := webcontext.PageInfo{Locale: "en", Dir: "ltr", Title: "Create your account", Path: onboardingPath}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	email := c.Query("email")
	first := c.Query("first")
	last := c.Query("last")
	method := c.Query("method")
	webutils.Render(c, http.StatusOK, pages.OnboardingStep1(email, first, last, method))
}

func (h *OnboardingHandler) Step2(c *gin.Context) {
	// Receive user fields and render step 2 with hidden inputs
	first := c.PostForm("first_name")
	last := c.PostForm("last_name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	method := c.PostForm("method")
	info := webcontext.PageInfo{Locale: "en", Dir: "ltr", Title: "Organization details", Path: onboardingPath}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, http.StatusOK, pages.OnboardingStep2(first, last, email, password, method))
}

func (h *OnboardingHandler) Complete(c *gin.Context) {
	// Gather all fields and create org+user atomically
	first := c.PostForm("first_name")
	last := c.PostForm("last_name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	method := c.PostForm("method")
	orgName := c.PostForm("org_name")
	slug := c.PostForm("org_slug")
	if first == "" || last == "" || email == "" || password == "" || orgName == "" || slug == "" {
		if method == "google" && password == "" && first != "" && last != "" && email != "" && orgName != "" && slug != "" {
			// generate a strong random password for oauth users
			if rp, err := utils.ID.RandomString(40); err == nil {
				password = rp
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
			return
		}
	}
	userReq := &account.CreateUserRequest{FirstName: first, LastName: last, Email: email, Password: password}
	orgReq := &account.CreateOrganizationRequest{Name: orgName, Slug: slug}
	user, err := h.onboarding.OnboardNewOrganization(c.Request.Context(), orgReq, userReq)
	if err != nil {
		if pd, ok := err.(*utils.ProblemDetails); ok {
			c.JSON(pd.Status, pd)
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "onboarding failed"})
		return
	}
	token, err2 := utils.JWT.GenerateToken(user.ID, user.OrganizationID)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	utils.JWT.SetJwtCookie(c, token)
	webutils.Redirect(c, "/")
}

func (h *OnboardingHandler) SlugAvailability(c *gin.Context) {
	slug := c.Query("slug")
	if slug == "" {
		c.String(http.StatusOK, "")
		return
	}
	ok, err := h.onboarding.IsOrganizationSlugAvailable(c.Request.Context(), slug)
	if err != nil {
		ok = false
	}
	if ok {
		c.String(http.StatusOK, "<div class=\"text-success text-sm\">Slug is available</div>")
	} else {
		c.String(http.StatusOK, "<div class=\"text-error text-sm\">Slug is taken</div>")
	}
}
