package e2e_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// OnboardingEmailVerifySuite tests POST /api/onboarding/email/verify endpoint
type OnboardingEmailVerifySuite struct {
	suite.Suite
	client *testutils.HTTPClient
	helper *OnboardingTestHelper
}

func (s *OnboardingEmailVerifySuite) SetupSuite() {
	s.client = testutils.NewHTTPClient("http://localhost:18080")
	s.helper = NewOnboardingTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *OnboardingEmailVerifySuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans")
}

func (s *OnboardingEmailVerifySuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans")
}

func (s *OnboardingEmailVerifySuite) TestVerifyEmail_Success() {
	otp := "123456"
	token, err := s.helper.CreateSessionWithOTP("verify@example.com", "starter", otp)
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
		"code":         otp,
		"password":     "SecurePassword123!",
		"firstName":    "John",
		"lastName":     "Doe",
	}
	resp, err := s.client.Post("/api/onboarding/email/verify", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	// Verify response structure
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal("identity_verified", result["stage"])
	s.Len(result, 1, "response should have exactly 1 field")
	s.Contains(result, "stage")

	// Verify session stage updated in database
	session, err := s.helper.GetSession(token)
	s.NoError(err)
	s.Equal("identity_verified", session["stage"])
}

func (s *OnboardingEmailVerifySuite) TestVerifyEmail_InvalidOTP() {
	correctOTP := "123456"
	token, err := s.helper.CreateSessionWithOTP("invalid@example.com", "starter", correctOTP)
	s.NoError(err)

	tests := []struct {
		name        string
		code        string
		expectedMsg string
	}{
		{"wrong code", "654321", "invalid code"},
		{"empty code", "", "code is required"},
		{"too short", "123", "code must be"},
		{"too long", "1234567", "code must be"},
		{"non-numeric", "abcdef", "invalid code"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload := map[string]interface{}{
				"sessionToken": token,
				"code":         tt.code,
				"password":     "SecurePassword123!",
				"firstName":    "John",
				"lastName":     "Doe",
			}
			resp, err := s.client.Post("/api/onboarding/email/verify", payload)
			s.NoError(err)
			defer resp.Body.Close()
			s.True(resp.StatusCode >= 400, "should reject invalid OTP: %s", tt.name)
		})
	}
}

func (s *OnboardingEmailVerifySuite) TestVerifyEmail_ExpiredOTP() {
	otp := "123456"
	token, err := s.helper.CreateSessionWithOTP("expired@example.com", "starter", otp)
	s.NoError(err)

	// Set OTP to expired
	err = s.helper.SetExpiredOTP(token)
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
		"code":         otp,
		"password":     "SecurePassword123!",
		"firstName":    "John",
		"lastName":     "Doe",
	}
	resp, err := s.client.Post("/api/onboarding/email/verify", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OnboardingEmailVerifySuite) TestVerifyEmail_WeakPasswords() {
	otp := "123456"
	tests := []struct {
		name       string
		password   string
		shouldFail bool
	}{
		{"too short", "Pass1!", true},           // Less than 8 chars
		{"empty", "", true},                     // Empty password
		{"no uppercase", "password123!", false}, // Valid - API only requires min 8 chars
		{"no lowercase", "PASSWORD123!", false}, // Valid - API only requires min 8 chars
		{"no number", "PasswordABC!", false},    // Valid - API only requires min 8 chars
		{"no special", "Password123", false},    // Valid - API only requires min 8 chars
		{"exactly 8 chars", "Pass123!", false},  // Valid - meets minimum
	}

	for i, tt := range tests {
		s.Run(tt.name, func() {
			email := fmt.Sprintf("weak_%d@example.com", i)
			token, err := s.helper.CreateSessionWithOTP(email, "starter", otp)
			s.NoError(err)

			payload := map[string]interface{}{
				"sessionToken": token,
				"code":         otp,
				"password":     tt.password,
				"firstName":    "John",
				"lastName":     "Doe",
			}
			resp, err := s.client.Post("/api/onboarding/email/verify", payload)
			s.NoError(err)
			defer resp.Body.Close()
			if tt.shouldFail {
				s.True(resp.StatusCode >= 400, "should reject weak password: %s", tt.name)
			} else {
				s.Equal(http.StatusOK, resp.StatusCode, "should accept password: %s", tt.name)
			}
		})
	}
}

func (s *OnboardingEmailVerifySuite) TestVerifyEmail_MissingFields() {
	otp := "123456"
	token, err := s.helper.CreateSessionWithOTP("missing@example.com", "starter", otp)
	s.NoError(err)

	tests := []struct {
		name    string
		payload map[string]interface{}
	}{
		{"missing code", map[string]interface{}{
			"sessionToken": token,
			"password":     "SecurePassword123!",
			"firstName":    "John",
			"lastName":     "Doe",
		}},
		{"missing password", map[string]interface{}{
			"sessionToken": token,
			"code":         otp,
			"firstName":    "John",
			"lastName":     "Doe",
		}},
		{"missing firstName", map[string]interface{}{
			"sessionToken": token,
			"code":         otp,
			"password":     "SecurePassword123!",
			"lastName":     "Doe",
		}},
		{"missing lastName", map[string]interface{}{
			"sessionToken": token,
			"code":         otp,
			"password":     "SecurePassword123!",
			"firstName":    "John",
		}},
		{"missing sessionToken", map[string]interface{}{
			"code":      otp,
			"password":  "SecurePassword123!",
			"firstName": "John",
			"lastName":  "Doe",
		}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.client.Post("/api/onboarding/email/verify", tt.payload)
			s.NoError(err)
			defer resp.Body.Close()
			s.Equal(http.StatusBadRequest, resp.StatusCode, "should require all fields: %s", tt.name)
		})
	}
}

func (s *OnboardingEmailVerifySuite) TestVerifyEmail_InvalidToken() {
	payload := map[string]interface{}{
		"sessionToken": "invalid_token_xyz",
		"code":         "123456",
		"password":     "SecurePassword123!",
		"firstName":    "John",
		"lastName":     "Doe",
	}
	resp, err := s.client.Post("/api/onboarding/email/verify", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *OnboardingEmailVerifySuite) TestVerifyEmail_InvalidStage() {
	// Create session at wrong stage (plan_selected instead of identity_pending)
	token, err := s.helper.CreateOnboardingSession("wrong@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
		"code":         "123456",
		"password":     "SecurePassword123!",
		"firstName":    "John",
		"lastName":     "Doe",
	}
	resp, err := s.client.Post("/api/onboarding/email/verify", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OnboardingEmailVerifySuite) TestVerifyEmail_XSSAttempts() {
	otp := "123456"
	token, err := s.helper.CreateSessionWithOTP("xss@example.com", "starter", otp)
	s.NoError(err)

	xssPayloads := []string{
		"<script>alert('xss')</script>",
		"John<img src=x onerror=alert('xss')>",
		"'; DROP TABLE users; --",
		"Doe</title><script>alert('xss')</script>",
	}

	for _, xss := range xssPayloads {
		payload := map[string]interface{}{
			"sessionToken": token,
			"code":         otp,
			"password":     "SecurePassword123!",
			"firstName":    xss,
			"lastName":     "Doe",
		}
		resp, err := s.client.Post("/api/onboarding/email/verify", payload)
		s.NoError(err)
		defer resp.Body.Close()

		// Should either reject or sanitize - check that it doesn't return 500
		s.NotEqual(http.StatusInternalServerError, resp.StatusCode, "should handle XSS safely")
	}
}

func (s *OnboardingEmailVerifySuite) TestVerifyEmail_NameValidation() {
	otp := "123456"
	token, err := s.helper.CreateSessionWithOTP("names@example.com", "starter", otp)
	s.NoError(err)

	tests := []struct {
		name       string
		firstName  string
		lastName   string
		shouldFail bool
	}{
		{"valid names", "John", "Doe", false},
		{"valid with spaces", "Mary Jane", "Van Der Berg", false},
		{"empty first name", "", "Doe", true},
		{"empty last name", "John", "", true},
		{"very long first name", string(make([]byte, 256)), "Doe", true},
		{"very long last name", "John", string(make([]byte, 256)), true},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload := map[string]interface{}{
				"sessionToken": token,
				"code":         otp,
				"password":     "SecurePassword123!",
				"firstName":    tt.firstName,
				"lastName":     tt.lastName,
			}
			resp, err := s.client.Post("/api/onboarding/email/verify", payload)
			s.NoError(err)
			defer resp.Body.Close()

			if tt.shouldFail {
				s.True(resp.StatusCode >= 400, "should reject: %s", tt.name)
			}
		})
	}
}

func (s *OnboardingEmailVerifySuite) TestVerifyEmail_ExpiredSession() {
	otp := "123456"
	token, err := s.helper.CreateSessionWithOTP("expired@example.com", "starter", otp)
	s.NoError(err)

	// Expire the session
	s.helper.ExpireSession(token)

	payload := map[string]interface{}{
		"sessionToken": token,
		"code":         otp,
		"password":     "SecurePassword123!",
		"firstName":    "John",
		"lastName":     "Doe",
	}
	resp, err := s.client.Post("/api/onboarding/email/verify", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestOnboardingEmailVerifySuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OnboardingEmailVerifySuite))
}
