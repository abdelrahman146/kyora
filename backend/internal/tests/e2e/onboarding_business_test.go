package e2e_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// OnboardingBusinessSuite tests POST /v1/onboarding/business endpoint
type OnboardingBusinessSuite struct {
	suite.Suite
	client *testutils.HTTPClient
	helper *OnboardingTestHelper
}

func (s *OnboardingBusinessSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient(e2eBaseURL)
	s.helper = NewOnboardingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *OnboardingBusinessSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingBusinessSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingBusinessSuite) TestSetBusiness_Success() {
	token, err := s.helper.CreateVerifiedSession("business@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
		"name":         "My Business",
		"descriptor":   "my-business",
		"country":      "AE",
		"currency":     "AED",
	}
	resp, err := s.client.Post("/v1/onboarding/business", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal("ready_to_commit", result["stage"], "free plan should go to ready_to_commit")
	// Verify response structure
	s.Len(result, 1, "response should have exactly 1 field")
	s.Contains(result, "stage")
}

func (s *OnboardingBusinessSuite) TestSetBusiness_PaidPlan() {
	// Create paid plan first
	s.helper.EnsureTestPlan("professional", "Professional", 54.99)

	token, err := s.helper.CreateVerifiedSession("paid@example.com", "professional")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
		"name":         "My Business",
		"descriptor":   "my-business",
		"country":      "AE",
		"currency":     "AED",
	}
	resp, err := s.client.Post("/v1/onboarding/business", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal("payment_pending", result["stage"], "paid plan should go to payment_pending")
	// Verify response structure
	s.Len(result, 1, "response should have exactly 1 field")
	s.Contains(result, "stage")
}

func (s *OnboardingBusinessSuite) TestSetBusiness_InvalidStage() {
	token, err := s.helper.CreateOnboardingSession("wrong@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
		"name":         "My Business",
		"descriptor":   "my-business",
		"country":      "AE",
		"currency":     "AED",
	}
	resp, err := s.client.Post("/v1/onboarding/business", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OnboardingBusinessSuite) TestSetBusiness_InvalidCountryCode() {
	tests := []struct {
		name       string
		country    string
		shouldFail bool
	}{
		{"too long", "USA", true},
		{"too short", "A", true},
		{"empty", "", true},
		{"lowercase", "ae", false}, // API accepts lowercase - only validates length
		{"numbers", "12", false},   // API accepts numbers - only validates length
		{"valid uppercase", "AE", false},
	}

	for i, tt := range tests {
		s.Run(tt.name, func() {
			// Create new session for each subtest to avoid stage conflicts
			token, err := s.helper.CreateVerifiedSession(fmt.Sprintf("country%d@example.com", i), "starter")
			s.NoError(err)

			payload := map[string]interface{}{
				"sessionToken": token,
				"name":         "My Business",
				"descriptor":   "my-business",
				"country":      tt.country,
				"currency":     "AED",
			}
			resp, err := s.client.Post("/v1/onboarding/business", payload)
			s.NoError(err)
			defer resp.Body.Close()
			if tt.shouldFail {
				s.Equal(http.StatusBadRequest, resp.StatusCode)
			} else {
				s.Equal(http.StatusOK, resp.StatusCode)
			}
		})
	}
}

func (s *OnboardingBusinessSuite) TestSetBusiness_InvalidCurrencyCode() {
	tests := []struct {
		name       string
		currency   string
		shouldFail bool
	}{
		{"too long", "AEDF", true},
		{"too short", "AE", true},
		{"empty", "", true},
		{"lowercase", "aed", false}, // API accepts lowercase - only validates length
		{"numbers", "123", false},   // API accepts numbers - only validates length
		{"valid uppercase", "AED", false},
	}

	for i, tt := range tests {
		s.Run(tt.name, func() {
			// Create new session for each subtest to avoid stage conflicts
			token, err := s.helper.CreateVerifiedSession(fmt.Sprintf("currency%d@example.com", i), "starter")
			s.NoError(err)

			payload := map[string]interface{}{
				"sessionToken": token,
				"name":         "My Business",
				"descriptor":   "my-business",
				"country":      "AE",
				"currency":     tt.currency,
			}
			resp, err := s.client.Post("/v1/onboarding/business", payload)
			s.NoError(err)
			defer resp.Body.Close()
			if tt.shouldFail {
				s.Equal(http.StatusBadRequest, resp.StatusCode)
			} else {
				s.Equal(http.StatusOK, resp.StatusCode)
			}
		})
	}
}

func (s *OnboardingBusinessSuite) TestSetBusiness_MissingFields() {
	token, err := s.helper.CreateVerifiedSession("missing@example.com", "starter")
	s.NoError(err)

	tests := []struct {
		name    string
		payload map[string]interface{}
	}{
		{
			"missing name",
			map[string]interface{}{
				"sessionToken": token,
				"descriptor":   "my-business",
				"country":      "AE",
				"currency":     "AED",
			},
		},
		{
			"missing descriptor",
			map[string]interface{}{
				"sessionToken": token,
				"name":         "My Business",
				"country":      "AE",
				"currency":     "AED",
			},
		},
		{
			"missing country",
			map[string]interface{}{
				"sessionToken": token,
				"name":         "My Business",
				"descriptor":   "my-business",
				"currency":     "AED",
			},
		},
		{
			"missing currency",
			map[string]interface{}{
				"sessionToken": token,
				"name":         "My Business",
				"descriptor":   "my-business",
				"country":      "AE",
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.client.Post("/v1/onboarding/business", tt.payload)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func (s *OnboardingBusinessSuite) TestSetBusiness_SQLInjection() {
	token, err := s.helper.CreateVerifiedSession("sql@example.com", "starter")
	s.NoError(err)

	maliciousInputs := []string{
		"'; DROP TABLE businesses; --",
		"business' OR '1'='1",
		"<script>alert('xss')</script>",
	}

	for _, malicious := range maliciousInputs {
		s.Run("injection_"+malicious[:10], func() {
			payload := map[string]interface{}{
				"sessionToken": token,
				"name":         malicious,
				"descriptor":   malicious,
				"country":      "AE",
				"currency":     "AED",
			}
			resp, err := s.client.Post("/v1/onboarding/business", payload)
			s.NoError(err)
			defer resp.Body.Close()
			// Should handle safely
			s.True(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest)
		})
	}
}

func TestOnboardingBusinessSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OnboardingBusinessSuite))
}
