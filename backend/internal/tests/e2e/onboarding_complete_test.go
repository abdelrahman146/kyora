package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// OnboardingCompleteSuite tests POST /v1/onboarding/complete endpoint
type OnboardingCompleteSuite struct {
	suite.Suite
	client *testutils.HTTPClient
	helper *OnboardingTestHelper
}

func (s *OnboardingCompleteSuite) SetupSuite() {
	s.client = testutils.NewHTTPClient(e2eBaseURL)
	s.helper = NewOnboardingTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *OnboardingCompleteSuite) SetupTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "sessions", "onboarding_sessions", "plans", "businesses", "shipping_zones", "subscriptions")
}

func (s *OnboardingCompleteSuite) TearDownTest() {
	testutils.TruncateTables(testEnv.Database, "users", "workspaces", "sessions", "onboarding_sessions", "plans", "businesses", "shipping_zones", "subscriptions")
}

func (s *OnboardingCompleteSuite) TestComplete_Success() {
	token, err := s.helper.CreateBusinessStagedSession("complete@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/v1/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.NotNil(result["user"], "should return user")
	s.NotEmpty(result["token"], "should return JWT token")
	s.NotEmpty(result["refreshToken"], "should return refresh token")
	// Verify response structure
	s.Len(result, 3, "response should have exactly 3 fields")
	s.Contains(result, "user")
	s.Contains(result, "token")
	s.Contains(result, "refreshToken")

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
	resp, err := s.client.Post("/v1/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	user := result["user"].(map[string]interface{})
	workspaceID := user["workspaceId"].(string)
	s.NotEmpty(workspaceID, "user should have workspaceId")

	count, err := s.helper.CountWorkspacesByID(workspaceID)
	s.NoError(err)
	s.Equal(int64(1), count, "workspace should exist in database")
}

func (s *OnboardingCompleteSuite) TestComplete_CreatesBusiness() {
	token, err := s.helper.CreateBusinessStagedSession("biz@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/v1/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	count, err := s.helper.CountBusinessesByName("Test Business")
	s.NoError(err)
	s.Equal(int64(1), count, "business should exist in database")
}

func (s *OnboardingCompleteSuite) TestComplete_CreatesDefaultShippingZone() {
	ctx := context.Background()

	token, err := s.helper.CreateBusinessStagedSession("zones@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/v1/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	user := result["user"].(map[string]interface{})
	workspaceID := user["workspaceId"].(string)
	s.NotEmpty(workspaceID)

	bizRepo := database.NewRepository[business.Business](testEnv.Database)
	biz, err := bizRepo.FindOne(ctx,
		bizRepo.ScopeEquals(business.BusinessSchema.WorkspaceID, workspaceID),
		bizRepo.ScopeEquals(business.BusinessSchema.Descriptor, "test-business"),
	)
	s.NoError(err)
	s.Equal("AE", biz.CountryCode)
	s.Equal("AED", biz.Currency)

	zoneRepo := database.NewRepository[business.ShippingZone](testEnv.Database)
	zones, err := zoneRepo.FindMany(ctx, zoneRepo.ScopeBusinessID(biz.ID))
	s.NoError(err)
	s.Len(zones, 1, "should create exactly one default shipping zone")
	s.Equal("AE", zones[0].Name)
	s.Equal([]string{"AE"}, []string(zones[0].Countries))
	s.Equal("AED", zones[0].Currency)
	s.True(zones[0].ShippingCost.IsZero())
	s.True(zones[0].FreeShippingThreshold.IsZero())
}

func (s *OnboardingCompleteSuite) TestComplete_InvalidToken() {
	payload := map[string]interface{}{
		"sessionToken": "invalid_token_12345",
	}
	resp, err := s.client.Post("/v1/onboarding/complete", payload)
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
	resp, err := s.client.Post("/v1/onboarding/complete", payload)
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
	resp, err := s.client.Post("/v1/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *OnboardingCompleteSuite) TestComplete_MissingToken() {
	payload := map[string]interface{}{}
	resp, err := s.client.Post("/v1/onboarding/complete", payload)
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
	resp, err := s.client.Post("/v1/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	jwtToken := result["token"].(string)
	s.NotEmpty(jwtToken, "should return JWT token")
	refreshToken := result["refreshToken"].(string)
	s.NotEmpty(refreshToken, "should return refresh token")

	// Token should be usable for authenticated requests
	s.Greater(len(jwtToken), 20, "JWT token should be reasonable length")

	// Verify response structure
	s.Len(result, 3, "response should have exactly 3 fields")
	s.Contains(result, "user")
	s.Contains(result, "token")
	s.Contains(result, "refreshToken")
}

func (s *OnboardingCompleteSuite) TestComplete_IdempotencySafety() {
	token, err := s.helper.CreateBusinessStagedSession("idempotent@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}

	// First completion should succeed
	resp1, err := s.client.Post("/v1/onboarding/complete", payload)
	s.NoError(err)
	defer resp1.Body.Close()
	s.Equal(http.StatusOK, resp1.StatusCode)

	// Second completion with same token should fail
	resp2, err := s.client.Post("/v1/onboarding/complete", payload)
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
	resp, err := s.client.Post("/v1/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	userCount, err := s.helper.CountUsersByEmail("consistency@example.com")
	s.NoError(err)
	s.Equal(int64(1), userCount, "exactly one user should be created")

	workspaceCount, err := s.helper.CountAllWorkspaces()
	s.NoError(err)
	s.Equal(int64(1), workspaceCount, "exactly one workspace should be created")

	businessCount, err := s.helper.CountAllBusinesses()
	s.NoError(err)
	s.Equal(int64(1), businessCount, "exactly one business should be created")
}

func (s *OnboardingCompleteSuite) TestComplete_SessionCleanedUp() {
	token, err := s.helper.CreateBusinessStagedSession("cleanup@example.com", "starter")
	s.NoError(err)

	payload := map[string]interface{}{
		"sessionToken": token,
	}
	resp, err := s.client.Post("/v1/onboarding/complete", payload)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	sess, err := s.helper.GetSessionModel(token)
	s.NoError(err)
	s.Equal("committed", string(sess.Stage))
	s.NotNil(sess.CommittedAt, "committed_at should be set")
}

func TestOnboardingCompleteSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(OnboardingCompleteSuite))
}
