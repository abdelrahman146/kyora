package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/abdelrahman146/kyora/internal/web/webutils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	// Google OAuth (scaffold)
	r.GET("/auth/google", h.GoogleAuth)
	r.GET("/auth/google/callback", h.GoogleCallback)
}

func (h *AuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	_, token, err := h.auth.Authenticate(c.Request.Context(), email, password)
	if err != nil {
		// return problem json
		if pd, ok := err.(*utils.ProblemDetails); ok {
			c.JSON(pd.Status, pd)
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
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
	// If org fields provided, complete immediately, else redirect to onboarding
	orgName := c.PostForm("org_name")
	slug := c.PostForm("org_slug")
	if orgName != "" && slug != "" {
		userReq := &account.CreateUserRequest{FirstName: first, LastName: last, Email: email, Password: password}
		orgReq := &account.CreateOrganizationRequest{Name: orgName, Slug: slug}
		user, err := h.onboard.OnboardNewOrganization(c.Request.Context(), orgReq, userReq)
		if err != nil {
			if pd, ok := err.(*utils.ProblemDetails); ok {
				c.JSON(pd.Status, pd)
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "registration failed"})
			return
		}
		token, _ := utils.JWT.GenerateToken(user.ID, user.OrganizationID)
		utils.JWT.SetJwtCookie(c, token)
		webutils.Redirect(c, "/")
		return
	}
	// No org provided: continue with onboarding wizard
	webutils.Redirect(c, "/onboarding")
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	email := c.PostForm("email")
	token, err := h.auth.CreateResetToken(c.Request.Context(), email)
	if err != nil {
		if pd, ok := err.(*utils.ProblemDetails); ok {
			c.JSON(pd.Status, pd)
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not initiate reset"})
		return
	}
	// Render a small success alert. In development, show the token to ease testing.
	dev := viper.GetString("env") != "production"
	msg := "If an account exists for that email, a reset link has been sent."
	if dev {
		msg = fmt.Sprintf("%s Dev token: %s", msg, token)
	}
	c.String(http.StatusOK, "<div class=\"alert alert-success mt-4\">%s</div>", msg)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	token := c.PostForm("token")
	newPassword := c.PostForm("password")
	confirm := c.PostForm("password_confirm")
	if newPassword == "" || newPassword != confirm {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}
	if err := h.auth.ConsumeResetToken(c.Request.Context(), token, newPassword); err != nil {
		if pd, ok := err.(*utils.ProblemDetails); ok {
			c.JSON(pd.Status, pd)
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "reset failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	cfg := h.googleConfig()
	if cfg == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Google OAuth not configured"})
		return
	}
	state := randomState()
	// store in cookie short-lived
	http.SetCookie(c.Writer, &http.Cookie{Name: "oauth_state", Value: state, Path: "/", HttpOnly: true, MaxAge: 300})
	authURL := cfg.AuthCodeURL(state, oauth2.AccessTypeOnline)
	c.Redirect(http.StatusFound, authURL)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	cfg := h.googleConfig()
	if cfg == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Google OAuth not configured"})
		return
	}
	// validate state
	state := c.Query("state")
	stateCookie, _ := c.Request.Cookie("oauth_state")
	if stateCookie == nil || stateCookie.Value == "" || stateCookie.Value != state {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}
	token, err := cfg.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code exchange failed"})
		return
	}
	client := cfg.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp.StatusCode != 200 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch userinfo"})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var info struct {
		Email      string `json:"email"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Name       string `json:"name"`
		Verified   bool   `json:"verified_email"`
	}
	_ = json.Unmarshal(body, &info)
	if info.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email not available"})
		return
	}
	// Try to find existing user
	user, err := h.auth.GetUserByEmail(c, info.Email)
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

func (h *AuthHandler) googleConfig() *oauth2.Config {
	clientID := viper.GetString("google.client_id")
	secret := viper.GetString("google.client_secret")
	redirect := viper.GetString("google.redirect_url")
	if clientID == "" || secret == "" || redirect == "" {
		return nil
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		RedirectURL:  redirect,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func randomState() string {
	s, _ := utils.ID.RandomString(24)
	return s
}
