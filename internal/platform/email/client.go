package email

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/spf13/viper"
)

// New returns a Client based on configuration.
// Priority:
// 1) email.mock.enabled = true -> Mock
// 2) email.provider = "mock" -> Mock
// 3) email.provider = "resend" (default) -> Resend
func New() (Client, error) {
	// allow both dotted and flat env keys to mirror existing codebase quirks
	mockEnabled := viper.GetBool(config.EmailMockEnabled)
	if !mockEnabled {
		// fallback to flat key
		if !mockEnabled {
			mockEnabled = viper.GetBool("email_mock_enabled")
		}
	}
	if mockEnabled {
		return &MockClient{}, nil
	}

	provider := viper.GetString(config.EmailProvider)
	if provider == "" {
		provider = viper.GetString("email_provider")
	}
	if provider == "" {
		provider = "resend"
	}

	switch provider {
	case "mock":
		return &MockClient{}, nil
	case "resend":
		apiKey := viper.GetString(config.ResendAPIKey)
		if apiKey == "" {
			apiKey = viper.GetString("email_resend_api_key")
		}
		if apiKey == "" {
			return nil, errors.New("missing Resend API key: set email.resend.api_key")
		}
		baseURL := viper.GetString(config.ResendAPIBaseURL)
		if baseURL == "" {
			baseURL = viper.GetString("email_resend_base_url")
		}
		if baseURL == "" {
			baseURL = "https://api.resend.com"
		}
		httpClient := &http.Client{Timeout: 15 * time.Second}
		return &ResendClient{apiKey: apiKey, baseURL: baseURL, httpClient: httpClient}, nil
	default:
		return nil, errors.New("unsupported email provider: " + provider)
	}
}

// Helper to expose a simple health check contract for providers
func Ping(ctx context.Context, c Client) error {
	// lightweight no-op for now; future: provider-specific status
	_, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	return nil
}
