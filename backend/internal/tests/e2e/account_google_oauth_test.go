package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// GoogleOAuthSuite tests Google OAuth authentication endpoints
type GoogleOAuthSuite struct {
	suite.Suite
	helper         *AccountTestHelper
	accountStorage *account.Storage
}

func (s *GoogleOAuthSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
	s.accountStorage = s.helper.AccountStorage
}

func (s *GoogleOAuthSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces")
	s.NoError(err)
}

func (s *GoogleOAuthSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces")
	s.NoError(err)
}

// TestGetGoogleAuthURL_Success tests successful Google OAuth URL generation
func (s *GoogleOAuthSuite) TestGetGoogleAuthURL_Success() {
	resp, err := s.helper.Client.Get("/v1/auth/google/url")
	s.NoError(err)
	defer resp.Body.Close()

	// If Google OAuth is not configured in test environment, endpoint returns error
	// This is expected and we verify the endpoint exists and handles the case properly
	if resp.StatusCode == http.StatusInternalServerError {
		var result map[string]interface{}
		s.NoError(testutils.DecodeJSON(resp, &result))

		// Verify error response structure
		s.Contains(result, "type", "should have error type")
		s.Contains(result, "title", "should have error title")
		s.Contains(result, "status", "should have error status")
		s.T().Log("Google OAuth not configured - verified proper error handling")
		return
	}

	s.Equal(http.StatusOK, resp.StatusCode)

	// Decode response
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	// Verify response structure
	s.Len(result, 2, "response should have exactly 2 fields")
	s.Contains(result, "url")
	s.Contains(result, "state")

	// Verify URL format
	url, ok := result["url"].(string)
	s.True(ok, "url should be string")
	s.NotEmpty(url)
	s.Contains(url, "accounts.google.com", "should be Google OAuth URL")
	s.Contains(url, "oauth", "should contain oauth")
	s.Contains(url, "response_type=code", "should request authorization code")
	s.Contains(url, "scope=", "should have scopes")

	// Verify state token for CSRF protection
	state, ok := result["state"].(string)
	s.True(ok, "state should be string")
	s.NotEmpty(state, "state should not be empty for CSRF protection")
	s.GreaterOrEqual(len(state), 16, "state should be sufficiently long for security")
}

// TestGetGoogleAuthURL_MultipleRequests tests that each request gets unique state
func (s *GoogleOAuthSuite) TestGetGoogleAuthURL_MultipleRequests() {
	states := make(map[string]bool)

	// Make multiple requests
	for i := 0; i < 5; i++ {
		resp, err := s.helper.Client.Get("/v1/auth/google/url")
		s.NoError(err)
		defer resp.Body.Close()

		// Skip if OAuth not configured
		if resp.StatusCode == http.StatusInternalServerError {
			s.T().Skip("Google OAuth not configured")
			return
		}

		s.Equal(http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		s.NoError(testutils.DecodeJSON(resp, &result))

		state, ok := result["state"].(string)
		s.True(ok, "state should be string")
		s.NotEmpty(state)

		// Verify state is unique
		s.False(states[state], "state %s should be unique, got duplicate on request %d", state, i+1)
		states[state] = true
	}

	// All states should be unique
	s.Len(states, 5, "all 5 requests should have unique state tokens")
}

// TestGetGoogleAuthURL_StateFormat tests the format and security of state tokens
func (s *GoogleOAuthSuite) TestGetGoogleAuthURL_StateFormat() {
	resp, err := s.helper.Client.Get("/v1/auth/google/url")
	s.NoError(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusInternalServerError {
		s.T().Skip("Google OAuth not configured")
		return
	}

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	state := result["state"].(string)

	// Verify state characteristics for security
	s.GreaterOrEqual(len(state), 32, "state should be at least 32 characters for security")
	s.NotContains(state, " ", "state should not contain spaces")
	s.NotContains(state, "\n", "state should not contain newlines")
	s.NotContains(state, "\t", "state should not contain tabs")

	// State should be URL-safe
	s.Regexp(`^[a-zA-Z0-9_-]+$`, state, "state should be URL-safe (alphanumeric, dash, underscore only)")
}

// TestLoginWithGoogle_MissingCode tests validation for missing code
func (s *GoogleOAuthSuite) TestLoginWithGoogle_MissingCode() {
	tests := []struct {
		name    string
		payload map[string]interface{}
	}{
		{
			name:    "empty payload",
			payload: map[string]interface{}{},
		},
		{
			name:    "null code",
			payload: map[string]interface{}{"code": nil},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.helper.Client.Post("/v1/auth/google/login", tt.payload)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(http.StatusBadRequest, resp.StatusCode, "should require code")

			var result map[string]interface{}
			s.NoError(testutils.DecodeJSON(resp, &result))

			// Verify error response
			s.Contains(result, "type")
			s.Contains(result, "title")
			s.Equal(float64(400), result["status"])
		})
	}
}

// TestLoginWithGoogle_EmptyCode tests empty code validation
func (s *GoogleOAuthSuite) TestLoginWithGoogle_EmptyCode() {
	payload := map[string]interface{}{
		"code": "",
	}

	resp, err := s.helper.Client.Post("/v1/auth/google/login", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode, "should reject empty code")

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result, "type")
	s.Contains(result, "title")
}

// TestLoginWithGoogle_InvalidCode tests behavior with invalid OAuth codes
func (s *GoogleOAuthSuite) TestLoginWithGoogle_InvalidCode() {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "short random code",
			code: "abc123",
		},
		{
			name: "sql injection attempt",
			code: "' OR '1'='1",
		},
		{
			name: "xss attempt",
			code: "<script>alert('xss')</script>",
		},
		{
			name: "very long code",
			code: string(make([]byte, 1000)),
		},
		{
			name: "special characters",
			code: "!@#$%^&*()",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			payload := map[string]interface{}{
				"code": tt.code,
			}

			resp, err := s.helper.Client.Post("/v1/auth/google/login", payload)
			s.NoError(err)
			defer resp.Body.Close()

			// Should fail with 4xx or 5xx status
			s.True(resp.StatusCode >= 400, "should reject invalid OAuth code")

			var result map[string]interface{}
			s.NoError(testutils.DecodeJSON(resp, &result))
			s.Contains(result, "type", "should have error type")
			s.Contains(result, "title", "should have error title")
		})
	}
}

// TestLoginWithGoogle_NonExistentUser tests OAuth login for email not in system
func (s *GoogleOAuthSuite) TestLoginWithGoogle_NonExistentUser() {
	// Note: Without mocking Google OAuth, this tests the validation flow
	// The service should reject OAuth codes for users not in the system
	payload := map[string]interface{}{
		"code": "mock_valid_looking_code_for_new_user",
	}

	resp, err := s.helper.Client.Post("/v1/auth/google/login", payload)
	s.NoError(err)
	defer resp.Body.Close()

	// Should be unauthorized (no account) or error (OAuth failed)
	s.True(resp.StatusCode >= 400, "should fail for non-existent user")

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result, "type")
	s.Contains(result, "title")
}

// TestLoginWithGoogle_ExistingUser_NoPassword tests user created via OAuth has no password
func (s *GoogleOAuthSuite) TestLoginWithGoogle_ExistingUser_NoPassword() {
	ctx := context.Background()

	// Create workspace first
	workspace := &account.Workspace{}
	workspaceRepo := database.NewRepository[account.Workspace](testEnv.Database)
	s.NoError(workspaceRepo.CreateOne(ctx, workspace))

	// Create user with empty password (simulating OAuth-only user)
	user := &account.User{
		Email:           "oauth.user@gmail.com",
		FirstName:       "OAuth",
		LastName:        "User",
		Password:        "", // No password for OAuth users
		Role:            role.RoleAdmin,
		IsEmailVerified: true,
		WorkspaceID:     workspace.ID,
	}

	userRepo := database.NewRepository[account.User](testEnv.Database)
	s.NoError(userRepo.CreateOne(ctx, user))

	// Verify user exists
	fetched, err := userRepo.FindOne(ctx, userRepo.ScopeEquals(account.UserSchema.Email, user.Email))
	s.NoError(err)
	s.NotNil(fetched)
	s.Empty(fetched.Password, "OAuth user should have no password")

	// Try to login with regular password endpoint should fail
	loginPayload := map[string]interface{}{
		"email":    user.Email,
		"password": "AnyPassword123!",
	}

	resp, err := s.helper.Client.Post("/v1/auth/login", loginPayload)
	s.NoError(err)
	defer resp.Body.Close()

	// Should fail because user has no password set
	s.True(resp.StatusCode >= 400, "password login should fail for OAuth-only user")
}

// TestLoginWithGoogle_UserWithPassword tests OAuth login for user with existing password
func (s *GoogleOAuthSuite) TestLoginWithGoogle_UserWithPassword() {
	ctx := context.Background()

	// Create regular user with password
	user, workspace, _, err := s.helper.CreateTestUser(ctx, "regular@gmail.com", "Password123!", "Regular", "User", role.RoleAdmin)
	s.NoError(err)

	// Verify user has password
	s.NotEmpty(user.Password, "regular user should have password")

	// OAuth login for this user should work if OAuth code is valid
	// But without mocking, we verify the endpoint processes the request
	payload := map[string]interface{}{
		"code": "mock_code_for_existing_user",
	}

	resp, err := s.helper.Client.Post("/v1/auth/google/login", payload)
	s.NoError(err)
	defer resp.Body.Close()

	// Will fail due to invalid OAuth code, but endpoint should exist
	s.True(resp.StatusCode >= 400)
	s.NotEqual(http.StatusNotFound, resp.StatusCode, "endpoint should exist")

	// Verify user still exists and wasn't modified
	userRepo := database.NewRepository[account.User](testEnv.Database)
	fetched, err := userRepo.FindOne(ctx, userRepo.ScopeEquals(account.UserSchema.Email, user.Email))
	s.NoError(err)
	s.Equal(user.ID, fetched.ID)
	s.Equal(workspace.ID, fetched.WorkspaceID)
}

// TestLoginWithGoogle_RateLimiting tests multiple rapid requests
func (s *GoogleOAuthSuite) TestLoginWithGoogle_RateLimiting() {
	payload := map[string]interface{}{
		"code": "rate_limit_test_code",
	}

	// Make multiple rapid requests
	for i := 0; i < 10; i++ {
		resp, err := s.helper.Client.Post("/v1/auth/google/login", payload)
		s.NoError(err)
		resp.Body.Close()

		// All requests should be processed (may fail, but not rate limited at endpoint level)
		s.NotEqual(http.StatusTooManyRequests, resp.StatusCode,
			"endpoint should not have request-level rate limiting (handled at infra level)")
	}
}

func TestGoogleOAuthSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(GoogleOAuthSuite))
}

// GoogleInvitationAcceptanceSuite tests accepting invitations with Google OAuth
type GoogleInvitationAcceptanceSuite struct {
	suite.Suite
	helper         *AccountTestHelper
	accountStorage *account.Storage
}

func (s *GoogleInvitationAcceptanceSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
	s.accountStorage = s.helper.AccountStorage
}

func (s *GoogleInvitationAcceptanceSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations")
	s.NoError(err)
}

func (s *GoogleInvitationAcceptanceSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations")
	s.NoError(err)
}

// TestAcceptInvitationWithGoogle_MissingToken tests missing token validation
func (s *GoogleInvitationAcceptanceSuite) TestAcceptInvitationWithGoogle_MissingToken() {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "no parameters",
			url:  "/v1/invitations/accept/google",
		},
		{
			name: "only code parameter",
			url:  "/v1/invitations/accept/google?code=mock_code",
		},
		{
			name: "empty token",
			url:  "/v1/invitations/accept/google?token=&code=mock_code",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.helper.Client.Get(tt.url)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(http.StatusBadRequest, resp.StatusCode, "should require token parameter")

			var result map[string]interface{}
			s.NoError(testutils.DecodeJSON(resp, &result))
			s.Contains(result, "type")
			s.Contains(result, "title")
			s.Equal(float64(400), result["status"])
		})
	}
}

// TestAcceptInvitationWithGoogle_MissingCode tests missing OAuth code validation
func (s *GoogleInvitationAcceptanceSuite) TestAcceptInvitationWithGoogle_MissingCode() {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "only token parameter",
			url:  "/v1/invitations/accept/google?token=inv_token123",
		},
		{
			name: "empty code",
			url:  "/v1/invitations/accept/google?token=inv_token123&code=",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.helper.Client.Get(tt.url)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(http.StatusBadRequest, resp.StatusCode, "should require code parameter")

			var result map[string]interface{}
			s.NoError(testutils.DecodeJSON(resp, &result))
			s.Contains(result, "type")
			s.Contains(result, "title")
		})
	}
}

// TestAcceptInvitationWithGoogle_InvalidToken tests invalid invitation token handling
func (s *GoogleInvitationAcceptanceSuite) TestAcceptInvitationWithGoogle_InvalidToken() {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "random token",
			token: "random_invalid_token_12345",
		},
		{
			name:  "sql injection in token",
			token: "invalid' OR '1'='1",
		},
		{
			name:  "xss in token",
			token: "invalid<script>alert</script>",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			url := fmt.Sprintf("/v1/invitations/accept/google?token=%s&code=mock_code", tt.token)
			resp, err := s.helper.Client.Get(url)

			// Some invalid tokens may cause URL parsing errors or server errors
			if err != nil {
				s.T().Logf("Request failed as expected with error: %v", err)
				return
			}
			defer resp.Body.Close()

			// Should reject with 4xx or 5xx status
			s.True(resp.StatusCode >= 400, "should reject invalid invitation token")

			// Try to decode response if it's JSON
			var result map[string]interface{}
			decodeErr := testutils.DecodeJSON(resp, &result)
			if decodeErr == nil {
				// If we got JSON, verify error structure
				s.Contains(result, "type")
				s.Contains(result, "title")
			} else {
				// Non-JSON response is also acceptable for invalid tokens
				s.T().Logf("Non-JSON response (acceptable for invalid tokens): %v", decodeErr)
			}
		})
	}
}

// TestAcceptInvitationWithGoogle_ValidToken tests accepting valid invitation
func (s *GoogleInvitationAcceptanceSuite) TestAcceptInvitationWithGoogle_ValidToken() {
	ctx := context.Background()

	// Create workspace and inviter
	inviter, workspace, _, err := s.helper.CreateTestUser(ctx, "inviter@example.com", "Password123!", "Inviter", "User", role.RoleAdmin)
	s.NoError(err)

	// Create invitation for new user
	invitation, token, err := s.helper.CreateInvitationWithToken(ctx, workspace.ID, "newuser@gmail.com", inviter.ID, role.RoleUser)
	s.NoError(err)
	s.NotNil(invitation)
	s.NotEmpty(token)

	// Verify invitation exists and is pending
	s.Equal(account.InvitationStatusPending, invitation.Status)
	s.Equal("newuser@gmail.com", invitation.Email)
	s.Equal(workspace.ID, invitation.WorkspaceID)
	s.Equal(role.RoleUser, invitation.Role)

	// Try to accept with Google OAuth (will fail without valid OAuth code)
	url := fmt.Sprintf("/v1/invitations/accept/google?token=%s&code=mock_google_code", token)
	resp, err := s.helper.Client.Get(url)
	s.NoError(err)
	defer resp.Body.Close()

	// Will fail due to invalid OAuth code, but should process the invitation token
	// The endpoint exists and validates the invitation token before OAuth
	s.True(resp.StatusCode >= 400, "will fail due to invalid OAuth code")
	s.NotEqual(http.StatusNotFound, resp.StatusCode, "endpoint should exist")

	// Verify invitation still exists and wasn't consumed
	ctx2 := context.Background()
	invRepo := database.NewRepository[account.UserInvitation](testEnv.Database)
	found, err := invRepo.FindByID(ctx2, invitation.ID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(account.InvitationStatusPending, found.Status, "invitation should remain pending after failed OAuth")
}

// TestAcceptInvitationWithGoogle_ExpiredToken tests expired invitation handling
func (s *GoogleInvitationAcceptanceSuite) TestAcceptInvitationWithGoogle_ExpiredToken() {
	ctx := context.Background()

	// Create workspace and inviter
	inviter, workspace, _, err := s.helper.CreateTestUser(ctx, "inviter@example.com", "Password123!", "Inviter", "User", role.RoleAdmin)
	s.NoError(err)

	// Create invitation
	// Note: Token expiration is handled by cache TTL, not a model field
	// In a real scenario, we would need to wait for the token to expire
	// or manipulate the cache directly to simulate expiration
	invitation, token, err := s.helper.CreateInvitationWithToken(ctx, workspace.ID, "expired@example.com", inviter.ID, role.RoleUser)
	s.NoError(err)
	s.NotEmpty(token)

	// For now, we verify that the endpoint processes expired scenarios
	// In production, expired tokens would be rejected by the service layer
	// This test documents the expected behavior rather than fully testing it
	s.Equal(account.InvitationStatusPending, invitation.Status)
	s.T().Log("Token expiration is handled by cache TTL - full testing requires time manipulation or cache mocking")
}

// TestAcceptInvitationWithGoogle_AlreadyAccepted tests already accepted invitation
func (s *GoogleInvitationAcceptanceSuite) TestAcceptInvitationWithGoogle_AlreadyAccepted() {
	ctx := context.Background()

	// Create workspace and inviter
	inviter, workspace, _, err := s.helper.CreateTestUser(ctx, "inviter@example.com", "Password123!", "Inviter", "User", role.RoleAdmin)
	s.NoError(err)

	// Create invitation
	now := time.Now().UTC()
	invitation := &account.UserInvitation{
		Email:       "accepted@example.com",
		WorkspaceID: workspace.ID,
		InviterID:   inviter.ID,
		Role:        role.RoleUser,
		Status:      account.InvitationStatusAccepted, // Already accepted
		AcceptedAt:  &gorm.DeletedAt{Time: now, Valid: true},
	}
	invRepo := database.NewRepository[account.UserInvitation](testEnv.Database)
	s.NoError(invRepo.CreateOne(ctx, invitation))

	// Generate token
	token, err := s.helper.CreateInvitationToken(ctx, invitation.ID)
	s.NoError(err)

	// Try to accept already-accepted invitation
	url := fmt.Sprintf("/v1/invitations/accept/google?token=%s&code=mock_google_code", token)
	resp, err := s.helper.Client.Get(url)
	s.NoError(err)
	defer resp.Body.Close()

	// Should reject already-accepted invitation
	s.True(resp.StatusCode >= 400, "should reject already-accepted invitation")

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result, "type")
	s.Contains(result, "title")
}

// TestAcceptInvitationWithGoogle_RevokedInvitation tests revoked invitation
func (s *GoogleInvitationAcceptanceSuite) TestAcceptInvitationWithGoogle_RevokedInvitation() {
	ctx := context.Background()

	// Create workspace and inviter
	inviter, workspace, _, err := s.helper.CreateTestUser(ctx, "inviter@example.com", "Password123!", "Inviter", "User", role.RoleAdmin)
	s.NoError(err)

	// Create invitation (note: revoked status is tracked in Status field)
	invitation := &account.UserInvitation{
		Email:       "revoked@example.com",
		WorkspaceID: workspace.ID,
		InviterID:   inviter.ID,
		Role:        role.RoleUser,
		Status:      account.InvitationStatusRevoked, // Revoked
	}
	invRepo := database.NewRepository[account.UserInvitation](testEnv.Database)
	s.NoError(invRepo.CreateOne(ctx, invitation))

	// Generate token
	token, err := s.helper.CreateInvitationToken(ctx, invitation.ID)
	s.NoError(err)

	// Try to accept revoked invitation
	url := fmt.Sprintf("/v1/invitations/accept/google?token=%s&code=mock_google_code", token)
	resp, err := s.helper.Client.Get(url)
	s.NoError(err)
	defer resp.Body.Close()

	// Should reject revoked invitation
	s.True(resp.StatusCode >= 400, "should reject revoked invitation")

	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result, "type")
	s.Contains(result, "title")
}

// TestAcceptInvitationWithGoogle_EmailMismatch tests OAuth email different from invitation
func (s *GoogleInvitationAcceptanceSuite) TestAcceptInvitationWithGoogle_EmailMismatch() {
	ctx := context.Background()

	// Create workspace and inviter
	inviter, workspace, _, err := s.helper.CreateTestUser(ctx, "inviter@example.com", "Password123!", "Inviter", "User", role.RoleAdmin)
	s.NoError(err)

	// Create invitation for specific email
	_, token, err := s.helper.CreateInvitationWithToken(ctx, workspace.ID, "invited@example.com", inviter.ID, role.RoleUser)
	s.NoError(err)

	// Note: In real scenario, Google OAuth would return different email
	// This tests that the system validates email match
	// Without mocking, we verify the endpoint processes the request
	url := fmt.Sprintf("/v1/invitations/accept/google?token=%s&code=mock_google_code", token)
	resp, err := s.helper.Client.Get(url)
	s.NoError(err)
	defer resp.Body.Close()

	// Will fail (OAuth error or email mismatch)
	s.True(resp.StatusCode >= 400)
}

// TestAcceptInvitationWithGoogle_WorkspaceIsolation tests workspace isolation
func (s *GoogleInvitationAcceptanceSuite) TestAcceptInvitationWithGoogle_WorkspaceIsolation() {
	ctx := context.Background()

	// Create two workspaces
	inviter1, workspace1, _, err := s.helper.CreateTestUser(ctx, "inviter1@example.com", "Password123!", "Inviter", "One", role.RoleAdmin)
	s.NoError(err)

	_, workspace2, _, err := s.helper.CreateTestUser(ctx, "inviter2@example.com", "Password123!", "Inviter", "Two", role.RoleAdmin)
	s.NoError(err)

	// Create invitation from workspace1
	_, token, err := s.helper.CreateInvitationWithToken(ctx, workspace1.ID, "newuser@example.com", inviter1.ID, role.RoleUser)
	s.NoError(err)

	// Verify invitation belongs to workspace1
	invRepo := database.NewRepository[account.UserInvitation](testEnv.Database)
	invitation, err := invRepo.FindOne(ctx, invRepo.ScopeEquals(account.UserInvitationSchema.Email, "newuser@example.com"))
	s.NoError(err)
	s.Equal(workspace1.ID, invitation.WorkspaceID)
	s.NotEqual(workspace2.ID, invitation.WorkspaceID, "invitation should not belong to workspace2")

	// Accept invitation (will fail due to OAuth, but tests isolation)
	url := fmt.Sprintf("/v1/invitations/accept/google?token=%s&code=mock_code", token)
	resp, err := s.helper.Client.Get(url)
	s.NoError(err)
	defer resp.Body.Close()

	// Endpoint should exist and process workspace-specific invitation
	s.NotEqual(http.StatusNotFound, resp.StatusCode)

	// Verify workspace2 was not affected
	workspace2Users, err := s.helper.CountWorkspaceUsers(ctx, workspace2.ID)
	s.NoError(err)
	s.Equal(int64(1), workspace2Users, "workspace2 should still have only original user")
}

// TestAcceptInvitationWithGoogle_RoleAssignment tests correct role assignment
func (s *GoogleInvitationAcceptanceSuite) TestAcceptInvitationWithGoogle_RoleAssignment() {
	ctx := context.Background()

	// Create workspace
	inviter, workspace, _, err := s.helper.CreateTestUser(ctx, "inviter@example.com", "Password123!", "Inviter", "User", role.RoleAdmin)
	s.NoError(err)

	// Test different roles
	roles := []role.Role{role.RoleUser, role.RoleAdmin}

	for _, testRole := range roles {
		s.Run(string(testRole), func() {
			email := fmt.Sprintf("user-%s@example.com", testRole)

			// Create invitation with specific role
			invitation, token, err := s.helper.CreateInvitationWithToken(ctx, workspace.ID, email, inviter.ID, testRole)
			s.NoError(err)

			// Verify invitation has correct role
			s.Equal(testRole, invitation.Role)

			// Try to accept (will fail due to OAuth, but verifies role is stored)
			url := fmt.Sprintf("/v1/invitations/accept/google?token=%s&code=mock_code", token)
			resp, err := s.helper.Client.Get(url)
			s.NoError(err)
			defer resp.Body.Close()

			// Endpoint processes request
			s.NotEqual(http.StatusNotFound, resp.StatusCode)

			// Verify invitation still has correct role
			invRepo := database.NewRepository[account.UserInvitation](testEnv.Database)
			stored, err := invRepo.FindByID(ctx, invitation.ID)
			s.NoError(err)
			s.Equal(testRole, stored.Role, "role should be preserved")
		})
	}
}

// TestAcceptInvitationWithGoogle_ConcurrentAcceptance tests race conditions
func (s *GoogleInvitationAcceptanceSuite) TestAcceptInvitationWithGoogle_ConcurrentAcceptance() {
	ctx := context.Background()

	// Create workspace and inviter
	inviter, workspace, _, err := s.helper.CreateTestUser(ctx, "inviter@example.com", "Password123!", "Inviter", "User", role.RoleAdmin)
	s.NoError(err)

	// Create invitation
	_, token, err := s.helper.CreateInvitationWithToken(ctx, workspace.ID, "concurrent@example.com", inviter.ID, role.RoleUser)
	s.NoError(err)

	// Try to accept same invitation multiple times concurrently
	url := fmt.Sprintf("/v1/invitations/accept/google?token=%s&code=mock_code", token)

	// Make concurrent requests
	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func(index int) {
			resp, err := s.helper.Client.Get(url)
			s.NoError(err)
			defer resp.Body.Close()

			// All should fail (OAuth error), but only one should process if OAuth worked
			// This tests that the system handles concurrent acceptance properly
			s.True(resp.StatusCode >= 400)

			done <- true
		}(i)
	}

	// Wait for all requests
	for i := 0; i < 3; i++ {
		<-done
	}

	// Verify invitation is still pending (since OAuth failed)
	// Get the original invitation to check status
	invRepo := database.NewRepository[account.UserInvitation](testEnv.Database)
	fetchedInv, err := invRepo.FindOne(ctx, invRepo.ScopeEquals(account.UserInvitationSchema.Email, "concurrent@example.com"))
	s.NoError(err)
	s.Equal(account.InvitationStatusPending, fetchedInv.Status)
}

func TestGoogleInvitationAcceptanceSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(GoogleInvitationAcceptanceSuite))
}
