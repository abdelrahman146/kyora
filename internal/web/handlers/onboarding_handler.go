package handlers

import (
	"net/http"
	"strings"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/abdelrahman146/kyora/internal/web/views/components"
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

type onboardingSubmission struct {
	First            string
	Last             string
	Email            string
	Password         string
	Method           string
	OrgName          string
	StoreName        string
	StoreCountryCode string
	StoreCurrency    string
}

func newOnboardingSubmission(c *gin.Context) onboardingSubmission {
	return onboardingSubmission{
		First:            strings.TrimSpace(c.PostForm("first_name")),
		Last:             strings.TrimSpace(c.PostForm("last_name")),
		Email:            strings.TrimSpace(c.PostForm("email")),
		Password:         c.PostForm("password"),
		Method:           c.PostForm("method"),
		OrgName:          strings.TrimSpace(c.PostForm("org_name")),
		StoreName:        strings.TrimSpace(c.PostForm("store_name")),
		StoreCountryCode: strings.TrimSpace(c.PostForm("store_country_code")),
		StoreCurrency:    strings.TrimSpace(c.PostForm("store_currency")),
	}
}

func (o *onboardingSubmission) normalize() {
	if o.StoreCurrency == "" {
		o.StoreCurrency = "USD"
	}
}

func (o *onboardingSubmission) validate() *utils.ProblemDetails {
	if o.First == "" || o.Email == "" || o.OrgName == "" || o.StoreName == "" || o.StoreCountryCode == "" {
		return utils.Problem.BadRequest("Please complete all required fields.")
	}
	if o.Method != "google" && o.Password == "" {
		return utils.Problem.BadRequest("Please add a password to secure your account.")
	}
	if o.Method == "google" && o.Password == "" {
		rp, err := utils.ID.RandomString(40)
		if err != nil {
			return utils.Problem.InternalError().WithError(err)
		}
		o.Password = rp
	}
	return nil
}

func (o onboardingSubmission) toRequests() (*account.CreateOrganizationRequest, *account.CreateUserRequest, *account.CreateInitialStoreRequest) {
	orgReq := &account.CreateOrganizationRequest{Name: o.OrgName}
	userReq := &account.CreateUserRequest{FirstName: o.First, LastName: o.Last, Email: o.Email, Password: o.Password}
	storeReq := &account.CreateInitialStoreRequest{
		Name:        o.StoreName,
		CountryCode: o.StoreCountryCode,
		Currency:    o.StoreCurrency,
	}
	return orgReq, userReq, storeReq
}

func NewOnboardingHandler(onboarding *account.OnboardingService, auth *account.AuthenticationService) *OnboardingHandler {
	return &OnboardingHandler{onboarding: onboarding, auth: auth}
}

func (h *OnboardingHandler) RegisterRoutes(r gin.IRoutes) {
	// wizard
	r.GET(onboardingPath, h.Index)
	r.POST(onboardingPath+"/step2", h.Step2)
	r.POST(onboardingPath+"/step3", h.Step3)
	r.POST(onboardingPath+"/complete", h.Complete)
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
	orgName := c.PostForm("org_name")
	info := webcontext.PageInfo{Locale: "en", Dir: "ltr", Title: "Organization details", Path: onboardingPath}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, http.StatusOK, pages.OnboardingStep2(first, last, email, password, method, orgName))
}

func (h *OnboardingHandler) Step3(c *gin.Context) {
	first := c.PostForm("first_name")
	last := c.PostForm("last_name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	method := c.PostForm("method")
	orgName := c.PostForm("org_name")
	slug := c.PostForm("org_slug")
	storeName := c.PostForm("store_name")
	storeCountryCode := c.PostForm("store_country_code")
	info := webcontext.PageInfo{Locale: "en", Dir: "ltr", Title: "Create your first store", Path: onboardingPath}
	ctx := webcontext.SetupPageInfo(c.Request.Context(), info)
	c.Request = c.Request.WithContext(ctx)
	webutils.Render(c, http.StatusOK, pages.OnboardingStep3(first, last, email, password, method, orgName, slug, storeName, storeCountryCode))
}

func (h *OnboardingHandler) Complete(c *gin.Context) {
	payload := newOnboardingSubmission(c)
	payload.normalize()
	if pd := payload.validate(); pd != nil {
		webutils.Render(c, pd.Status, components.Alert("error", pd.Detail))
		return
	}
	orgReq, userReq, storeReq := payload.toRequests()
	user, err := h.onboarding.OnboardNewOrganization(c.Request.Context(), orgReq, userReq, storeReq)
	if err != nil {
		if pd, ok := err.(*utils.ProblemDetails); ok {
			webutils.Render(c, pd.Status, components.Alert("error", pd.Detail))
			return
		}
		webutils.Render(c, http.StatusBadRequest, components.Alert("error", "Onboarding failed. Please try again."))
		return
	}
	token, err2 := utils.JWT.GenerateToken(user.ID, user.OrganizationID)
	if err2 != nil {
		webutils.Render(c, http.StatusInternalServerError, components.Alert("error", "Failed to complete onboarding."))
		return
	}
	utils.JWT.SetJwtCookie(c, token)
	webutils.Redirect(c, "/")
}
