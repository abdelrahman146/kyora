package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/abdelrahman146/kyora/internal/web/views/components"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type AuthHandler struct {
	auth    *account.AuthenticationService
	onboard *account.OnboardingService
}

func NewAuthHandler(auth *account.AuthenticationService, onboard *account.OnboardingService) *AuthHandler {
	return &AuthHandler{auth: auth, onboard: onboard}
}

func (h *AuthHandler) RegisterRoutes(r gin.IRoutes) {
	r.POST("/login", h.Login)
	r.POST("/register", h.Register)
	r.POST("/forgot-password", h.ForgotPassword)
	r.POST("/reset-password", h.ResetPassword)
	r.GET("/auth/google", h.GoogleAuth)
	r.GET("/auth/google/callback", h.GoogleCallback)
}

func (h *AuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	_, token, err := h.auth.Authenticate(c.Request.Context(), email, password)
	if err != nil {
		if pd, ok := err.(*utils.ProblemDetails); ok {
			webutils.RenderFragments(c, pd.Status, components.Alert("error", pd.Detail), components.AlertFragmentKey)
			return
		}
		webutils.RenderFragments(c, http.StatusUnauthorized, components.Alert("error", "invalid credentials"), components.AlertFragmentKey)
		return
	}
	utils.JWT.SetJwtCookie(c, token)
	webutils.Redirect(c, "/")
}

func (h *AuthHandler) Register(c *gin.Context) {
	first := c.PostForm("first_name")
	last := c.PostForm("last_name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	// Optional immediate org creation if provided
	orgName := c.PostForm("org_name")
	slug := c.PostForm("org_slug")
	if orgName != "" && slug != "" {
		userReq := &account.CreateUserRequest{FirstName: first, LastName: last, Email: email, Password: password}
		orgReq := &account.CreateOrganizationRequest{Name: orgName, Slug: slug}
		user, err := h.onboard.OnboardNewOrganization(c.Request.Context(), orgReq, userReq)
		if err != nil {
			if pd, ok := err.(*utils.ProblemDetails); ok {
				webutils.RenderFragments(c, pd.Status, components.Alert("error", pd.Detail), components.AlertFragmentKey)
				return
			}
			webutils.RenderFragments(c, http.StatusBadRequest, components.Alert("error", "registration failed"), components.AlertFragmentKey)
			return
		}
		token, _ := utils.JWT.GenerateToken(user.ID, user.OrganizationID)
		utils.JWT.SetJwtCookie(c, token)
		webutils.Redirect(c, "/")
		return
	}
	// Otherwise, continue onboarding wizard with provided user fields
	q := url.Values{}
	q.Set("first", first)
	q.Set("last", last)
	q.Set("email", email)
	webutils.Redirect(c, "/onboarding?"+q.Encode())
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	email := c.PostForm("email")
	token, err := h.auth.CreateResetToken(c.Request.Context(), email)
	if err != nil {
		if pd, ok := err.(*utils.ProblemDetails); ok {
			webutils.RenderFragments(c, pd.Status, components.Alert("error", pd.Detail), components.AlertFragmentKey)
			return
		}
		webutils.RenderFragments(c, http.StatusBadRequest, components.Alert("error", "could not initiate reset"), components.AlertFragmentKey)
		return
	}
	msg := "If the email exists, a reset token has been generated."
	if viper.GetString("env") != "production" {
		msg = fmt.Sprintf("%s Dev token: %s", msg, token)
	}
	webutils.RenderFragments(c, http.StatusOK, components.Alert("success", msg), components.AlertFragmentKey)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	token := c.PostForm("token")
	pwd := c.PostForm("password")
	confirm := c.PostForm("password_confirm")
	if pwd == "" || confirm == "" || pwd != confirm {
		webutils.RenderFragments(c, http.StatusBadRequest, components.Alert("error", "passwords do not match"), components.AlertFragmentKey)
		return
	}
	if err := h.auth.ConsumeResetToken(c.Request.Context(), token, pwd); err != nil {
		if pd, ok := err.(*utils.ProblemDetails); ok {
			webutils.RenderFragments(c, pd.Status, components.Alert("error", pd.Detail), components.AlertFragmentKey)
			return
		}
		webutils.RenderFragments(c, http.StatusBadRequest, components.Alert("error", "reset failed"), components.AlertFragmentKey)
		return
	}
	// small delay to mitigate token brute force timing
	time.Sleep(150 * time.Millisecond)
	webutils.RenderFragments(c, http.StatusOK, components.Alert("success", "password updated"), components.AlertFragmentKey)
}

func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	state, _ := utils.ID.RandomString(24)
	// Store state in a short-lived cookie
	http.SetCookie(c.Writer, &http.Cookie{Name: "oauth_state", Value: state, Path: "/", HttpOnly: true, MaxAge: 300, SameSite: http.SameSiteLaxMode})
	authURL, pd := h.auth.GoogleGetAuthURL(c.Request.Context(), state)
	if pd != nil {
		webutils.RenderFragments(c, pd.Status, components.Alert("error", pd.Detail), components.AlertFragmentKey)
		return
	}
	c.Redirect(http.StatusFound, authURL)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	if code == "" || state == "" {
		webutils.RenderFragments(c, http.StatusBadRequest, components.Alert("error", "missing code or state"), components.AlertFragmentKey)
		return
	}
	// verify state
	if cookie, _ := c.Request.Cookie("oauth_state"); cookie == nil || cookie.Value != state {
		webutils.RenderFragments(c, http.StatusBadRequest, components.Alert("error", "invalid state"), components.AlertFragmentKey)
		return
	}
	info, pd := h.auth.GoogleExchangeAndFetchUser(c.Request.Context(), code)
	if pd != nil {
		webutils.RenderFragments(c, pd.Status, components.Alert("error", pd.Detail), components.AlertFragmentKey)
		return
	}
	// Try to log in if user exists and has org
	user, err := h.auth.GetUserByEmail(c.Request.Context(), info.Email)
	if err == nil && user != nil && user.OrganizationID != "" {
		jwt, _ := utils.JWT.GenerateToken(user.ID, user.OrganizationID)
		utils.JWT.SetJwtCookie(c, jwt)
		webutils.Redirect(c, "/")
		return
	}
	// New user or missing org: send to onboarding with prefilled details
	first, last := info.GivenName, info.FamilyName
	if first == "" && last == "" && info.Name != "" {
		parts := strings.SplitN(info.Name, " ", 2)
		first = parts[0]
		if len(parts) > 1 {
			last = parts[1]
		}
	}
	q := url.Values{}
	q.Set("email", info.Email)
	q.Set("first", first)
	q.Set("last", last)
	q.Set("method", "google")
	c.Redirect(http.StatusFound, "/onboarding?"+q.Encode())
}
