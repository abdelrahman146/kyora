package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Name       string `json:"name"`
	Verified   bool   `json:"verified_email"`
}

func googleConfig() *oauth2.Config {
	clientID := viper.GetString(config.GoogleOAuthClientID)
	secret := viper.GetString(config.GoogleOAuthClientSecret)
	redirect := viper.GetString(config.GoogleOAuthRedirectURL)
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

func GoogleGetAuthURL(ctx context.Context, state string) (string, error) {
	cfg := googleConfig()
	if cfg == nil {
		return "", fmt.Errorf("google OAuth not configured")
	}
	return cfg.AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

func GoogleExchangeAndFetchUser(ctx context.Context, code string) (*GoogleUserInfo, error) {
	cfg := googleConfig()
	if cfg == nil {
		return nil, fmt.Errorf("google OAuth not configured")
	}
	tok, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}
	client := cfg.Client(ctx, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch userinfo: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userinfo error: %d %s", resp.StatusCode, string(b))
	}
	var info GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo: %w", err)
	}
	if info.Email == "" {
		return nil, fmt.Errorf("email not available from provider")
	}
	return &info, nil
}
