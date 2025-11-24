package e2e_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/tests/e2e"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// OnboardingEmailOTPSuite tests POST /api/onboarding/email/otp endpoint
type OnboardingEmailOTPSuite struct {
	suite.Suite
	client *testutils.HTTPClient
	helper *e2e.OnboardingTestHelper
}

func (s *OnboardingEmailOTPSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient("http://localhost:18080")
	s.helper = e2e.NewOnboardingTestHelper(testEnv.Database, "http://localhost:18080")
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
	resp, err := s.client.Post("/api/onboarding/email/otp", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func (s *OnboardingEmailOTPSuite) TestSendEmailOTP_MissingToken() {
	payload := map[string]interface{}{}
	resp, err := s.client.Post("/api/onboarding/email/otp", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OnboardingEmailOTPSuite) TestSendEmailOTP_InvalidToken() {
	payload := map[string]interface{}{
		"sessionToken": "invalid_token_12345",
	}
	resp, err := s.client.Post("/api/onboarding/email/otp", payload)
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
	resp, err := s.client.Post("/api/onboarding/email/otp", payload)
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
	resp1, err := s.client.Post("/api/onboarding/email/otp", payload)
	s.NoError(err)
	defer resp1.Body.Close()
	s.Equal(http.StatusNoContent, resp1.StatusCode)

	// Second request immediately after should be rate limited
	resp2, err := s.client.Post("/api/onboarding/email/otp", payload)
	s.NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusTooManyRequests, resp2.StatusCode)

	// Wait 1 second and retry - should still be rate limited (30s minimum)
	time.Sleep(1 * time.Second)
	resp3, err := s.client.Post("/api/onboarding/email/otp", payload)
	s.NoError(err)
	defer resp3.Body.Close()
	s.Equal(http.StatusTooManyRequests, resp3.StatusCode)
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
	resp, err := s.client.Post("/api/onboarding/email/otp", payload)
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
