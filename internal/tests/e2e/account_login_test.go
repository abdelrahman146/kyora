package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// LoginSuite tests the POST /v1/auth/login endpoint
type LoginSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *LoginSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *LoginSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces")
	s.NoError(err)
}

func (s *LoginSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces")
	s.NoError(err)
}

func (s *LoginSuite) TestLogin_Success() {
	ctx := context.Background()

	// Create test user
	_, _, _, err := s.helper.CreateTestUser(ctx, "test@example.com", "ValidPassword123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	// Attempt login
	payload := map[string]interface{}{
		"email":    "test@example.com",
		"password": "ValidPassword123!",
	}

	resp, err := s.helper.Client.Post("/v1/auth/login", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	// Decode and verify response
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	// Verify exact response structure
	s.Len(result, 2, "response should have exactly 2 fields")
	s.Contains(result, "user")
	s.Contains(result, "token")

	// Verify user object
	user := result["user"].(map[string]interface{})
	s.Equal("test@example.com", user["email"])
	s.Equal("John", user["firstName"])
	s.Equal("Doe", user["lastName"])
	s.NotEmpty(user["id"])
	s.NotEmpty(user["workspaceId"])
	s.NotEmpty(user["role"])

	// Verify JWT token is set in cookies
	s.NotEmpty(result["token"])
}

func (s *LoginSuite) TestLogin_InvalidEmail() {
	payload := map[string]interface{}{
		"email":    "nonexistent@example.com",
		"password": "ValidPassword123!",
	}

	resp, err := s.helper.Client.Post("/v1/auth/login", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *LoginSuite) TestLogin_InvalidPassword() {
	ctx := context.Background()

	// Create test user
	_, _, _, err := s.helper.CreateTestUser(ctx, "test@example.com", "ValidPassword123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	// Attempt login with wrong password
	payload := map[string]interface{}{
		"email":    "test@example.com",
		"password": "WrongPassword123!",
	}

	resp, err := s.helper.Client.Post("/v1/auth/login", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *LoginSuite) TestLogin_MissingEmail() {
	payload := map[string]interface{}{
		"password": "ValidPassword123!",
	}

	resp, err := s.helper.Client.Post("/v1/auth/login", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *LoginSuite) TestLogin_MissingPassword() {
	payload := map[string]interface{}{
		"email": "test@example.com",
	}

	resp, err := s.helper.Client.Post("/v1/auth/login", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *LoginSuite) TestLogin_InvalidEmailFormat() {
	payload := map[string]interface{}{
		"email":    "not-an-email",
		"password": "ValidPassword123!",
	}

	resp, err := s.helper.Client.Post("/v1/auth/login", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *LoginSuite) TestLogin_EmailEnumerationPrevention() {
	ctx := context.Background()

	// Create test user
	_, _, _, err := s.helper.CreateTestUser(ctx, "exists@example.com", "ValidPassword123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	// Try with invalid email
	payload1 := map[string]interface{}{
		"email":    "nonexistent@example.com",
		"password": "ValidPassword123!",
	}

	resp1, err := s.helper.Client.Post("/v1/auth/login", payload1)
	s.NoError(err)
	defer resp1.Body.Close()

	// Try with valid email but wrong password
	payload2 := map[string]interface{}{
		"email":    "exists@example.com",
		"password": "WrongPassword123!",
	}

	resp2, err := s.helper.Client.Post("/v1/auth/login", payload2)
	s.NoError(err)
	defer resp2.Body.Close()

	// Both should return same error to prevent email enumeration
	s.Equal(resp1.StatusCode, resp2.StatusCode)
	s.Equal(http.StatusUnauthorized, resp1.StatusCode)
}

func (s *LoginSuite) TestLogin_MultipleUsers_DifferentWorkspaces() {
	ctx := context.Background()

	// Create two users in different workspaces
	user1, workspace1, _, err := s.helper.CreateTestUser(ctx, "user1@example.com", "Password123!", "User", "One", role.RoleAdmin)
	s.NoError(err)

	user2, workspace2, _, err := s.helper.CreateTestUser(ctx, "user2@example.com", "Password123!", "User", "Two", role.RoleAdmin)
	s.NoError(err)

	// Verify they're in different workspaces
	s.NotEqual(workspace1.ID, workspace2.ID)
	s.NotEqual(user1.ID, user2.ID)

	// Login as user1
	payload1 := map[string]interface{}{
		"email":    "user1@example.com",
		"password": "Password123!",
	}

	resp1, err := s.helper.Client.Post("/v1/auth/login", payload1)
	s.NoError(err)
	defer resp1.Body.Close()

	s.Equal(http.StatusOK, resp1.StatusCode)
	var result1 map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp1, &result1))

	user1Data := result1["user"].(map[string]interface{})
	s.Equal(workspace1.ID, user1Data["workspaceId"])
}

func TestLoginSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(LoginSuite))
}
