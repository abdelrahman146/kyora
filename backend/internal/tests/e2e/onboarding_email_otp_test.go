package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// OnboardingEmailOTPSuite tests POST /v1/onboarding/email/otp endpoint
type OnboardingEmailOTPSuite struct {
	suite.Suite
	client *testutils.HTTPClient
	helper *OnboardingTestHelper
}

func (s *OnboardingEmailOTPSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient(e2eBaseURL)
	s.helper = NewOnboardingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *OnboardingEmailOTPSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingEmailOTPSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingEmailOTPSuite) TestSendEmailOTP_Success() {
	token, err := s.helper.CreateOnboardingSession("otp@example.com", "starter")
	s.NoError(err)
	s.NotEmpty(token)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/v1/onboarding/email/otp", payload)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	var body struct {
		RetryAfterSeconds int `json:"retryAfterSeconds"`
	}
	err = testutils.DecodeJSON(resp, &body)
	s.NoError(err)
	s.GreaterOrEqual(body.RetryAfterSeconds, 1)
}

func (s *OnboardingEmailOTPSuite) TestSendEmailOTP_MissingToken() {
	payload := map[string]interface{}{}
	resp, err := s.client.Post("/v1/onboarding/email/otp", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OnboardingEmailOTPSuite) TestSendEmailOTP_InvalidToken() {
	payload := map[string]interface{}{
		"sessionToken": "invalid_token_12345",
	}
	resp, err := s.client.Post("/v1/onboarding/email/otp", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *OnboardingEmailOTPSuite) TestSendEmailOTP_InvalidStage() {
	token, err := s.helper.CreateOnboardingSession("stage@example.com", "starter")
	s.NoError(err)
	s.NotEmpty(token)

	// Set to wrong stage
	s.helper.UpdateSessionStage(token, "identity_verified")

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/v1/onboarding/email/otp", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OnboardingEmailOTPSuite) TestSendEmailOTP_RateLimit() {
	token, err := s.helper.CreateOnboardingSession("rate@example.com", "starter")
	s.NoError(err)
	s.NotEmpty(token)

	payload := map[string]interface{}{
		"sessionToken": token,
	}

	// First request should succeed
	resp1, err := s.client.Post("/v1/onboarding/email/otp", payload)
	s.NoError(err)
	s.Equal(http.StatusOK, resp1.StatusCode)
	_ = resp1.Body.Close()

	// Second request immediately after should be rate limited
	resp2, err := s.client.Post("/v1/onboarding/email/otp", payload)
	s.NoError(err)
	s.Equal(http.StatusTooManyRequests, resp2.StatusCode)
	var p2 struct {
		Extensions map[string]any `json:"extensions"`
	}
	decErr := json.NewDecoder(resp2.Body).Decode(&p2)
	_ = resp2.Body.Close()
	s.NoError(decErr)
	seconds, ok := p2.Extensions["retryAfterSeconds"].(float64)
	s.True(ok)
	s.GreaterOrEqual(seconds, float64(0))
	s.LessOrEqual(seconds, float64(120))

	// Wait 1 second and retry - should still be rate limited (2m cooldown)
	time.Sleep(1 * time.Second)
	resp3, err := s.client.Post("/v1/onboarding/email/otp", payload)
	s.NoError(err)
	s.Equal(http.StatusTooManyRequests, resp3.StatusCode)
	_ = resp3.Body.Close()
}

func (s *OnboardingEmailOTPSuite) TestSendEmailOTP_ExpiredSession() {
	token, err := s.helper.CreateOnboardingSession("expired@example.com", "starter")
	s.NoError(err)
	s.NotEmpty(token)

	// Expire the session
	s.helper.ExpireSession(token)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/v1/onboarding/email/otp", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestOnboardingEmailOTPSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OnboardingEmailOTPSuite))
}
