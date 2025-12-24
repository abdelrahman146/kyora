package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// UserProfileSuite tests user profile endpoints
// GET /v1/users/me and PATCH /v1/users/me
type UserProfileSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *UserProfileSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *UserProfileSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces")
	s.NoError(err)
}

func (s *UserProfileSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces")
	s.NoError(err)
}

func (s *UserProfileSuite) TestGetProfile_Success() {
	ctx := context.Background()

	user, workspace, token, err := s.helper.CreateTestUser(ctx, "test@example.com", "Password123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/users/me", nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(user.ID, result["id"])
	s.Equal(user.Email, result["email"])
	s.Equal(user.FirstName, result["firstName"])
	s.Equal(user.LastName, result["lastName"])
	s.Equal(workspace.ID, result["workspaceId"])
	s.Equal("admin", result["role"])
	s.Equal(true, result["isEmailVerified"])
}

func (s *UserProfileSuite) TestGetProfile_Unauthenticated() {
	resp, err := s.helper.Client.Get("/v1/users/me")
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *UserProfileSuite) TestUpdateProfile_Success() {
	ctx := context.Background()

	user, _, token, err := s.helper.CreateTestUser(ctx, "test@example.com", "Password123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	payload := map[string]interface{}{
		"firstName": "Jane",
		"lastName":  "Smith",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/users/me", payload, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(user.ID, result["id"])
	s.Equal("Jane", result["firstName"])
	s.Equal("Smith", result["lastName"])
	s.Equal(user.Email, result["email"]) // Email should not change
}

func (s *UserProfileSuite) TestUpdateProfile_PartialUpdate() {
	ctx := context.Background()

	user, _, token, err := s.helper.CreateTestUser(ctx, "test@example.com", "Password123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	// Update only first name
	payload := map[string]interface{}{
		"firstName": "Jane",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/users/me", payload, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal("Jane", result["firstName"])
	s.Equal("Doe", result["lastName"]) // Last name should remain unchanged
	_ = user
}

func (s *UserProfileSuite) TestUpdateProfile_Unauthenticated() {
	payload := map[string]interface{}{
		"firstName": "Jane",
	}

	resp, err := s.helper.Client.Patch("/v1/users/me", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *UserProfileSuite) TestUpdateProfile_EmptyPayload() {
	ctx := context.Background()

	user, _, token, err := s.helper.CreateTestUser(ctx, "test@example.com", "Password123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	payload := map[string]interface{}{}

	resp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/users/me", payload, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	// Nothing should change
	s.Equal("John", result["firstName"])
	s.Equal("Doe", result["lastName"])
	_ = user
}

func TestUserProfileSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(UserProfileSuite))
}

// WorkspaceSuite tests workspace endpoints
// GET /v1/workspaces/me and GET /v1/workspaces/users
type WorkspaceSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *WorkspaceSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *WorkspaceSuite) SetupTest() {
	s.NoError(testEnv.Cache.FlushAll())
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces")
	s.NoError(err)
}

func (s *WorkspaceSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces")
	s.NoError(err)
}

func (s *WorkspaceSuite) TestGetWorkspace_Success() {
	ctx := context.Background()

	user, workspace, token, err := s.helper.CreateTestUser(ctx, "test@example.com", "Password123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/workspaces/me", nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(workspace.ID, result["id"])
	s.Equal(user.ID, result["ownerId"])
}

func (s *WorkspaceSuite) TestGetWorkspace_Unauthenticated() {
	resp, err := s.helper.Client.Get("/v1/workspaces/me")
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *WorkspaceSuite) TestGetWorkspaceUsers_Success() {
	ctx := context.Background()

	workspace, users, err := testutils.CreateWorkspaceWithUsers(ctx, testEnv.Database, "owner@example.com", "Password123!", []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Role      role.Role
	}{
		{"user1@example.com", "Password123!", "User", "One", role.RoleUser},
		{"admin1@example.com", "Password123!", "Admin", "One", role.RoleAdmin},
	})
	s.NoError(err)

	// Login as owner
	token, err := testutils.LoginAndGetToken(s.helper.Client, "owner@example.com", "Password123!")
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/workspaces/users", nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var result []map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	// Should have 3 users (owner + 2 created)
	s.Len(result, 3)

	// Verify all users belong to same workspace
	for _, u := range result {
		s.Equal(workspace.ID, u["workspaceId"])
	}

	_ = users
}

func (s *WorkspaceSuite) TestGetWorkspaceUsers_Unauthenticated() {
	resp, err := s.helper.Client.Get("/v1/workspaces/users")
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *WorkspaceSuite) TestGetWorkspaceUser_Success() {
	ctx := context.Background()

	workspace, users, err := testutils.CreateWorkspaceWithUsers(ctx, testEnv.Database, "owner@example.com", "Password123!", []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Role      role.Role
	}{
		{"user1@example.com", "Password123!", "User", "One", role.RoleUser},
	})
	s.NoError(err)

	// Login as owner
	token, err := testutils.LoginAndGetToken(s.helper.Client, "owner@example.com", "Password123!")
	s.NoError(err)

	targetUser := users[1]

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/workspaces/users/"+targetUser.ID, nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	s.Equal(targetUser.ID, result["id"])
	s.Equal(targetUser.Email, result["email"])
	s.Equal(workspace.ID, result["workspaceId"])
}

func (s *WorkspaceSuite) TestGetWorkspaceUser_NotFound() {
	ctx := context.Background()

	_, _, token, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/workspaces/users/nonexistent-id", nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *WorkspaceSuite) TestGetWorkspaceUser_DifferentWorkspace() {
	ctx := context.Background()

	// Create two separate workspaces
	_, _, token1, err := s.helper.CreateTestUser(ctx, "user1@example.com", "Password123!", "User", "One", role.RoleAdmin)
	s.NoError(err)

	user2, _, _, err := s.helper.CreateTestUser(ctx, "user2@example.com", "Password123!", "User", "Two", role.RoleAdmin)
	s.NoError(err)

	// Try to access user2 from workspace1
	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/workspaces/users/"+user2.ID, nil, token1)
	s.NoError(err)
	defer resp.Body.Close()

	// Should not find user from different workspace
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func TestWorkspaceSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(WorkspaceSuite))
}
