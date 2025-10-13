package account

import (
	"context"
	"encoding/json"
	"time"

	"io"
	"net/http"

	"github.com/abdelrahman146/kyora/internal/db"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthenticationService struct {
	userRepo *userRepository
	cache    *db.Memcache
}

func NewAuthenticationService(userRepo *userRepository, cache *db.Memcache) *AuthenticationService {
	return &AuthenticationService{userRepo: userRepo, cache: cache}
}

func (s *AuthenticationService) Authenticate(ctx context.Context, email, password string) (*User, string, error) {
	user, err := s.userRepo.findOne(ctx, s.userRepo.scopeEmail(email), db.WithPreload(OrganizationStruct))
	if err != nil {
		return nil, "", err
	}
	if user == nil || !utils.Hash.Validate(password, user.PasswordHash) {
		return nil, "", utils.Problem.Unauthorized("Invalid email or password")
	}
	jwt, err := utils.JWT.GenerateToken(user.ID, user.OrganizationID)
	if err != nil {
		return nil, "", err
	}
	return user, jwt, nil
}

func (s *AuthenticationService) GetUserByID(ctx context.Context, id string) (*User, error) {
	user, err := s.userRepo.findByID(ctx, id, db.WithPreload(OrganizationStruct))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthenticationService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.userRepo.findOne(ctx, s.userRepo.scopeEmail(email), db.WithPreload(OrganizationStruct))
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Google OAuth
type GoogleUserInfo struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Name       string `json:"name"`
	Verified   bool   `json:"verified_email"`
}

func (s *AuthenticationService) googleConfig() *oauth2.Config {
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

func (s *AuthenticationService) GoogleGetAuthURL(ctx context.Context, state string) (string, *utils.ProblemDetails) {
	cfg := s.googleConfig()
	if cfg == nil {
		return "", utils.Problem.BadRequest("Google OAuth not configured")
	}
	return cfg.AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

func (s *AuthenticationService) GoogleExchangeAndFetchUser(ctx context.Context, code string) (*GoogleUserInfo, *utils.ProblemDetails) {
	cfg := s.googleConfig()
	if cfg == nil {
		return nil, utils.Problem.BadRequest("Google OAuth not configured")
	}
	tok, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, utils.Problem.BadRequest("code exchange failed").WithError(err)
	}
	client := cfg.Client(ctx, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, utils.Problem.BadRequest("failed to fetch userinfo").WithError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, utils.Problem.BadRequest("userinfo error").With("status", resp.StatusCode).With("body", string(b))
	}
	var info GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, utils.Problem.InternalError().WithError(err)
	}
	if info.Email == "" {
		return nil, utils.Problem.BadRequest("email not available from provider")
	}
	return &info, nil
}

type resetPayload struct {
	UserID         string    `json:"userId"`
	OrganizationID string    `json:"organizationId"`
	Email          string    `json:"email"`
	ExpAt          time.Time `json:"expAt"`
}

const resetPrefix = "pwreset:"

func (s *AuthenticationService) CreateResetToken(ctx context.Context, email string) string {
	user, err := s.userRepo.findOne(ctx, s.userRepo.scopeEmail(email))
	if err != nil {
		return ""
	}
	token, err := utils.ID.RandomString(40)
	if err != nil {
		utils.Log.FromContext(ctx).Error("failed to generate reset token", "err", err)
		return ""
	}
	passwordResetTtlSeconds := viper.GetInt64("password_reset_ttl_seconds")
	ttl := time.Duration(passwordResetTtlSeconds) * time.Second
	payload := resetPayload{UserID: user.ID, OrganizationID: user.OrganizationID, Email: user.Email, ExpAt: time.Now().Add(ttl)}
	b, _ := json.Marshal(payload)
	exp := int32(ttl.Seconds())
	if exp <= 0 {
		exp = 900 // default 15m
	}
	if err := s.cache.Set(resetPrefix+token, b, exp); err != nil {
		utils.Log.FromContext(ctx).Error("failed to store reset token", "err", err)
		return ""
	}
	return token
}

func (s *AuthenticationService) ValidateResetToken(ctx context.Context, token string) (*User, error) {
	b, err := s.cache.Get(resetPrefix + token)
	if err != nil || len(b) == 0 {
		return nil, utils.Problem.NotFound("Invalid or expired token")
	}
	var p resetPayload
	_ = json.Unmarshal(b, &p)
	if time.Now().After(p.ExpAt) {
		_ = s.cache.Delete(resetPrefix + token)
		return nil, utils.Problem.NotFound("Token expired")
	}
	user, err := s.userRepo.findByID(ctx, p.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthenticationService) ConsumeResetToken(ctx context.Context, token string, newPassword string) error {
	user, err := s.ValidateResetToken(ctx, token)
	if err != nil {
		return err
	}
	hash, err := utils.Hash.Make(newPassword)
	if err != nil {
		return utils.Problem.InternalError().WithError(err)
	}
	updates := &User{PasswordHash: hash}
	if err := s.userRepo.patchOne(ctx, updates, s.userRepo.scopeID(user.ID)); err != nil {
		return err
	}
	_ = s.cache.Delete(resetPrefix + token)
	return nil
}
