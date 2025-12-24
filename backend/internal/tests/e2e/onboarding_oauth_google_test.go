package e2e_test

import (
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// OnboardingOAuthGoogleSuite tests POST /api/onboarding/oauth/google endpoint.
// Note: Test environment might not have Google OAuth configured; we focus on validation and session gating.
type OnboardingOAuthGoogleSuite struct {
	suite.Suite
	client *testutils.HTTPClient
	helper *OnboardingTestHelper
}

func (s *OnboardingOAuthGoogleSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient(e2eBaseURL)
	s.helper = NewOnboardingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *OnboardingOAuthGoogleSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingOAuthGoogleSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses")
}

func (s *OnboardingOAuthGoogleSuite) TestOAuthGoogle_ValidationErrors() {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name:           "missing payload fields",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing sessionToken",
			payload:        map[string]interface{}{"code": "abc"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing code",
			payload:        map[string]interface{}{"sessionToken": "tok"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty code",
			payload:        map[string]interface{}{"sessionToken": "tok", "code": ""},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.client.Post("/api/onboarding/oauth/google", tt.payload)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(tt.expectedStatus, resp.StatusCode)

			var result map[string]interface{}
			s.NoError(testutils.DecodeJSON(resp, &result))
			s.Contains(result, "type")
			s.Contains(result, "title")
			s.Contains(result, "status")
			s.Equal(float64(tt.expectedStatus), result["status"])
		})
	}
}

func (s *OnboardingOAuthGoogleSuite) TestOAuthGoogle_SessionNotFound() {
	resp, err := s.client.Post("/api/onboarding/oauth/google", map[string]interface{}{
		"sessionToken": "sess_does_not_exist",
		"code":         "any",
	})
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result, "type")
	s.Contains(result, "title")
	s.Contains(result, "status")
	s.Equal(float64(http.StatusNotFound), result["status"])
}

func (s *OnboardingOAuthGoogleSuite) TestOAuthGoogle_InvalidCode_Rejected() {
	// Create a real session token first
	token, err := s.helper.CreateOnboardingSession("oauthuser@example.com", "starter")
	s.NoError(err)
	s.NotEmpty(token)

	resp, err := s.client.Post("/api/onboarding/oauth/google", map[string]interface{}{
		"sessionToken": token,
		"code":         "<script>alert('xss')</script>",
	})
	s.NoError(err)
	defer resp.Body.Close()

	// Without mocking Google OAuth, this should fail with an error (4xx/5xx) but must not succeed.
	s.True(resp.StatusCode >= 400)
}

func TestOnboardingOAuthGoogleSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OnboardingOAuthGoogleSuite))
}
