package e2e_test

import (
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// OnboardingCompleteSuite tests POST /api/onboarding/complete endpoint
type OnboardingCompleteSuite struct {
	suite.Suite
	client *testutils.HTTPClient
	helper *OnboardingTestHelper
}

func (s *OnboardingCompleteSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient("http://localhost:18080")
	s.helper = NewOnboardingTestHelper(testEnv.Database, "http://localhost:18080")
}

func (s *OnboardingCompleteSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses", "subscriptions")
}

func (s *OnboardingCompleteSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "onboarding_sessions", "plans", "businesses", "subscriptions")
}

func (s *OnboardingCompleteSuite) TestComplete_Success() {
	token, err := s.helper.CreateBusinessStagedSession("complete@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.NotNil(result["user"], "should return user")
	s.NotEmpty(result["token"], "should return JWT token")
	// Verify response structure
	s.Len(result, 2, "response should have exactly 2 fields")
	s.Contains(result, "user")
	s.Contains(result, "token")

	user := result["user"].(map[string]interface{})
	s.Equal("complete@example.com", user["email"])
	s.Equal("Test", user["firstName"])
	s.Equal("User", user["lastName"])
	s.Equal(true, user["isEmailVerified"])
	// Verify all user fields are present
	s.NotEmpty(user["id"], "user should have id")
	s.NotEmpty(user["workspaceId"], "user should have workspaceId")
	s.NotEmpty(user["role"], "user should have role")
	s.Contains(user, "email")
	s.Contains(user, "firstName")
	s.Contains(user, "lastName")
	s.Contains(user, "isEmailVerified")
}

func (s *OnboardingCompleteSuite) TestComplete_CreatesWorkspace() {
	token, err := s.helper.CreateBusinessStagedSession("workspace@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	user := result["user"].(map[string]interface{})
	workspaceID := user["workspaceId"].(string)
	s.NotEmpty(workspaceID, "user should have workspaceId")

	// Verify workspace exists in database
	db := testEnv.Database.GetDB()
	var count int64
	db.Raw("SELECT COUNT(*) FROM workspaces WHERE id = ?", workspaceID).Scan(&count)
	s.Equal(int64(1), count, "workspace should exist in database")
}

func (s *OnboardingCompleteSuite) TestComplete_CreatesBusiness() {
	token, err := s.helper.CreateBusinessStagedSession("biz@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	// Verify business exists in database
	db := testEnv.Database.GetDB()
	var count int64
	db.Raw("SELECT COUNT(*) FROM businesses WHERE name = ?", "Test Business").Scan(&count)
	s.Equal(int64(1), count, "business should exist in database")
}

func (s *OnboardingCompleteSuite) TestComplete_InvalidToken() {
	payload := map[string]interface{}{
		"sessionToken": "invalid_token_12345",
	}
	resp, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *OnboardingCompleteSuite) TestComplete_WrongStage() {
	token, err := s.helper.CreateOnboardingSession("wrong@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OnboardingCompleteSuite) TestComplete_ExpiredSession() {
	token, err := s.helper.CreateBusinessStagedSession("expired@example.com", "starter")
	s.NoError(err)

	// Expire the session
	s.helper.ExpireSession(token)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OnboardingCompleteSuite) TestComplete_MissingToken() {
	payload := map[string]interface{}{}
	resp, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OnboardingCompleteSuite) TestComplete_ValidJWTToken() {
	token, err := s.helper.CreateBusinessStagedSession("jwt@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	jwtToken := result["token"].(string)
	s.NotEmpty(jwtToken, "should return JWT token")

	// Token should be usable for authenticated requests
	s.Greater(len(jwtToken), 20, "JWT token should be reasonable length")

	// Verify response structure
	s.Len(result, 2, "response should have exactly 2 fields")
	s.Contains(result, "user")
	s.Contains(result, "token")
}

func (s *OnboardingCompleteSuite) TestComplete_IdempotencySafety() {
	token, err := s.helper.CreateBusinessStagedSession("idempotent@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}

	// First completion should succeed
	resp1, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp1.Body.Close()
	s.Equal(http.StatusOK, resp1.StatusCode)

	// Second completion with same token should fail
	resp2, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp2.Body.Close()
	s.True(resp2.StatusCode >= 400, "should not allow duplicate completion")
}

func (s *OnboardingCompleteSuite) TestComplete_DatabaseConsistency() {
	token, err := s.helper.CreateBusinessStagedSession("consistency@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	// Verify all related records were created properly
	db := testEnv.Database.GetDB()

	var userCount, workspaceCount, businessCount int64
	db.Raw("SELECT COUNT(*) FROM users WHERE email = ?", "consistency@example.com").Scan(&userCount)
	s.Equal(int64(1), userCount, "exactly one user should be created")

	db.Raw("SELECT COUNT(*) FROM workspaces").Scan(&workspaceCount)
	s.Equal(int64(1), workspaceCount, "exactly one workspace should be created")

	db.Raw("SELECT COUNT(*) FROM businesses").Scan(&businessCount)
	s.Equal(int64(1), businessCount, "exactly one business should be created")
}

func (s *OnboardingCompleteSuite) TestComplete_SessionCleanedUp() {
	token, err := s.helper.CreateBusinessStagedSession("cleanup@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/api/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	// Verify session is marked as committed
	db := testEnv.Database.GetDB()
	var stage string
	var committedAt *string
	db.Raw("SELECT stage, committed_at FROM onboarding_sessions WHERE token = ?", token).Row().Scan(&stage, &committedAt)
	s.Equal("committed", stage)
	s.NotNil(committedAt, "committed_at should be set")
}

func TestOnboardingCompleteSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OnboardingCompleteSuite))
}
