package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// RefreshTokenSuite tests refresh token rotation and logout endpoints.
//
// Internally, refresh tokens are persisted as Sessions (hashed token records).
// These tests focus on the external contract and security behavior.
type RefreshTokenSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *RefreshTokenSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, e2eBaseURL)
}

func (s *RefreshTokenSuite) SetupTest() {
	s.NoError(testEnv.Cache.FlushAll())
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "sessions")
	s.NoError(err)
}

func (s *RefreshTokenSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "sessions")
	s.NoError(err)
}

func (s *RefreshTokenSuite) TestRefresh_RotatesRefreshToken() {
	ctx := context.Background()

	// Create test user
	_, _, _, err := s.helper.CreateTestUser(ctx, "rt@example.com", "ValidPassword123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	// Login
	loginResp, err := s.helper.Client.Post("/v1/auth/login", map[string]interface{}{
		"email":    "rt@example.com",
		"password": "ValidPassword123!",
	})
	s.NoError(err)
	defer loginResp.Body.Close()
	s.Equal(http.StatusOK, loginResp.StatusCode)

	var loginResult map[string]interface{}
	s.NoError(testutils.DecodeJSON(loginResp, &loginResult))
	s.Len(loginResult, 3, "response should have exactly 3 fields")
	accessToken := loginResult["token"].(string)
	refreshToken := loginResult["refreshToken"].(string)
	s.NotEmpty(accessToken)
	s.NotEmpty(refreshToken)

	// Refresh (rotates refresh token)
	refreshResp, err := s.helper.Client.Post("/v1/auth/refresh", map[string]interface{}{"refreshToken": refreshToken})
	s.NoError(err)
	defer refreshResp.Body.Close()
	s.Equal(http.StatusOK, refreshResp.StatusCode)

	var refreshed map[string]interface{}
	s.NoError(testutils.DecodeJSON(refreshResp, &refreshed))
	s.Len(refreshed, 2, "response should have exactly 2 fields")

	newAccess := refreshed["token"].(string)
	newRefresh := refreshed["refreshToken"].(string)
	s.NotEmpty(newAccess)
	s.NotEmpty(newRefresh)
	s.NotEqual(accessToken, newAccess)
	s.NotEqual(refreshToken, newRefresh)

	// Old refresh token must be invalid after rotation
	refreshResp2, err := s.helper.Client.Post("/v1/auth/refresh", map[string]interface{}{"refreshToken": refreshToken})
	s.NoError(err)
	defer refreshResp2.Body.Close()
	s.Equal(http.StatusUnauthorized, refreshResp2.StatusCode)
}

func (s *RefreshTokenSuite) TestLogout_RevokesRefreshToken() {
	ctx := context.Background()

	_, _, _, err := s.helper.CreateTestUser(ctx, "logout@example.com", "ValidPassword123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	loginResp, err := s.helper.Client.Post("/v1/auth/login", map[string]interface{}{
		"email":    "logout@example.com",
		"password": "ValidPassword123!",
	})
	s.NoError(err)
	defer loginResp.Body.Close()
	s.Equal(http.StatusOK, loginResp.StatusCode)

	var loginResult map[string]interface{}
	s.NoError(testutils.DecodeJSON(loginResp, &loginResult))
	refreshToken := loginResult["refreshToken"].(string)
	s.NotEmpty(refreshToken)

	logoutResp, err := s.helper.Client.Post("/v1/auth/logout", map[string]interface{}{"refreshToken": refreshToken})
	s.NoError(err)
	defer logoutResp.Body.Close()
	s.Equal(http.StatusNoContent, logoutResp.StatusCode)

	refreshResp, err := s.helper.Client.Post("/v1/auth/refresh", map[string]interface{}{"refreshToken": refreshToken})
	s.NoError(err)
	defer refreshResp.Body.Close()
	s.Equal(http.StatusUnauthorized, refreshResp.StatusCode)
}

func (s *RefreshTokenSuite) TestLogoutOtherDevices_RevokesOthersOnly() {
	ctx := context.Background()

	_, _, _, err := s.helper.CreateTestUser(ctx, "others@example.com", "ValidPassword123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	// Device A login
	loginA, err := s.helper.Client.Post("/v1/auth/login", map[string]interface{}{
		"email":    "others@example.com",
		"password": "ValidPassword123!",
	})
	s.NoError(err)
	defer loginA.Body.Close()
	s.Equal(http.StatusOK, loginA.StatusCode)
	var resA map[string]interface{}
	s.NoError(testutils.DecodeJSON(loginA, &resA))
	rtA := resA["refreshToken"].(string)
	s.NotEmpty(rtA)

	// Device B login
	loginB, err := s.helper.Client.Post("/v1/auth/login", map[string]interface{}{
		"email":    "others@example.com",
		"password": "ValidPassword123!",
	})
	s.NoError(err)
	defer loginB.Body.Close()
	s.Equal(http.StatusOK, loginB.StatusCode)
	var resB map[string]interface{}
	s.NoError(testutils.DecodeJSON(loginB, &resB))
	rtB := resB["refreshToken"].(string)
	s.NotEmpty(rtB)

	logoutOthers, err := s.helper.Client.Post("/v1/auth/logout-others", map[string]interface{}{"refreshToken": rtA})
	s.NoError(err)
	defer logoutOthers.Body.Close()
	s.Equal(http.StatusNoContent, logoutOthers.StatusCode)

	// Device A still can refresh
	refreshA, err := s.helper.Client.Post("/v1/auth/refresh", map[string]interface{}{"refreshToken": rtA})
	s.NoError(err)
	defer refreshA.Body.Close()
	s.Equal(http.StatusOK, refreshA.StatusCode)

	// Device B can no longer refresh
	refreshB, err := s.helper.Client.Post("/v1/auth/refresh", map[string]interface{}{"refreshToken": rtB})
	s.NoError(err)
	defer refreshB.Body.Close()
	s.Equal(http.StatusUnauthorized, refreshB.StatusCode)
}

func TestRefreshTokenSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(RefreshTokenSuite))
}
