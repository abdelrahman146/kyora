package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// InvitationManagementSuite tests invitation management endpoints
// POST /v1/workspaces/invitations, GET /v1/workspaces/invitations, and DELETE /v1/workspaces/invitations/:id
type InvitationManagementSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *InvitationManagementSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *InvitationManagementSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations", "subscriptions", "plans")
	s.NoError(err)
}

func (s *InvitationManagementSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations", "subscriptions", "plans")
	s.NoError(err)
}

func (s *InvitationManagementSuite) TestInviteUser_Success() {
	ctx := context.Background()
	_, workspace, token, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	// Create subscription for workspace
	err = s.helper.CreateTestSubscription(ctx, workspace.ID)
	s.NoError(err)

	payload := map[string]interface{}{
		"email": "newuser@example.com",
		"role":  "user",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/workspaces/invitations", payload, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusCreated, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	// Assert invitation fields
	s.Equal("newuser@example.com", result["email"])
	s.Equal("user", result["role"])
	s.Equal("pending", result["status"])
	s.NotEmpty(result["id"])
}

func (s *InvitationManagementSuite) TestInviteUser_DuplicateEmail() {
	ctx := context.Background()
	user, workspace, token, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	// Create subscription for workspace
	err = s.helper.CreateTestSubscription(ctx, workspace.ID)
	s.NoError(err)

	// First invitation should succeed
	payload := map[string]interface{}{
		"email": "newuser@example.com",
		"role":  "user",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/workspaces/invitations", payload, token)
	s.NoError(err)
	s.Equal(http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// Try to invite the same email again - should return conflict
	resp, err = s.helper.Client.AuthenticatedRequest("POST", "/v1/workspaces/invitations", payload, token)
	s.NoError(err)
	defer resp.Body.Close()

	// Should return conflict because a pending invitation already exists for this email
	s.Equal(http.StatusConflict, resp.StatusCode)
	_ = user
	_ = workspace
}

func (s *InvitationManagementSuite) TestInviteUser_NoPermission() {
	ctx := context.Background()

	// Create regular user without manage permission
	_, _, token, err := s.helper.CreateTestUser(ctx, "user@example.com", "Password123!", "Regular", "User", role.RoleUser)
	s.NoError(err)

	payload := map[string]interface{}{
		"email": "newuser@example.com",
		"role":  "user",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("POST", "/v1/workspaces/invitations", payload, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusForbidden, resp.StatusCode)
}

func (s *InvitationManagementSuite) TestGetWorkspaceInvitations_Success() {
	ctx := context.Background()
	user, workspace, token, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	// Create subscription for workspace
	err = s.helper.CreateTestSubscription(ctx, workspace.ID)
	s.NoError(err)

	// Create some invitations directly in database
	_, err = s.helper.CreateInvitation(ctx, workspace.ID, "invite1@example.com", user.ID, role.RoleUser, account.InvitationStatusPending)
	s.NoError(err)
	_, err = s.helper.CreateInvitation(ctx, workspace.ID, "invite2@example.com", user.ID, role.RoleAdmin, account.InvitationStatusPending)
	s.NoError(err)
	_, err = s.helper.CreateInvitation(ctx, workspace.ID, "invite3@example.com", user.ID, role.RoleUser, account.InvitationStatusAccepted)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/workspaces/invitations", nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)
	var result []map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	// Should return all invitations (3)
	s.Len(result, 3)
}

func (s *InvitationManagementSuite) TestGetWorkspaceInvitations_FilterByStatus() {
	ctx := context.Background()
	user, workspace, token, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	// Create subscription for workspace
	err = s.helper.CreateTestSubscription(ctx, workspace.ID)
	s.NoError(err)

	// Create invitations with different statuses
	_, err = s.helper.CreateInvitation(ctx, workspace.ID, "pending1@example.com", user.ID, role.RoleUser, account.InvitationStatusPending)
	s.NoError(err)
	_, err = s.helper.CreateInvitation(ctx, workspace.ID, "pending2@example.com", user.ID, role.RoleUser, account.InvitationStatusPending)
	s.NoError(err)
	_, err = s.helper.CreateInvitation(ctx, workspace.ID, "accepted@example.com", user.ID, role.RoleUser, account.InvitationStatusAccepted)
	s.NoError(err)

	// Filter for pending invitations only
	resp, err := s.helper.Client.AuthenticatedRequest("GET", "/v1/workspaces/invitations?status=pending", nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)
	var result []map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	// Should return only pending invitations (2)
	s.Len(result, 2)
	for _, inv := range result {
		s.Equal("pending", inv["status"])
	}
}

func (s *InvitationManagementSuite) TestRevokeInvitation_Success() {
	ctx := context.Background()
	user, workspace, token, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	// Create invitation
	invitation, err := s.helper.CreateInvitation(ctx, workspace.ID, "invite@example.com", user.ID, role.RoleUser, account.InvitationStatusPending)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/workspaces/invitations/"+invitation.ID, nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify invitation status is revoked
	updatedInvitation, err := s.helper.GetInvitation(ctx, invitation.ID)
	s.NoError(err)
	s.Equal(account.InvitationStatusRevoked, updatedInvitation.Status)
}

func (s *InvitationManagementSuite) TestRevokeInvitation_NotFound() {
	ctx := context.Background()
	_, _, token, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/workspaces/invitations/nonexistent-id", nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *InvitationManagementSuite) TestRevokeInvitation_AlreadyAccepted() {
	ctx := context.Background()
	user, workspace, token, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	// Create accepted invitation
	invitation, err := s.helper.CreateInvitation(ctx, workspace.ID, "invite@example.com", user.ID, role.RoleUser, account.InvitationStatusAccepted)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/workspaces/invitations/"+invitation.ID, nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusConflict, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result["detail"], "only pending invitations can be revoked")
}

func TestInvitationManagementSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(InvitationManagementSuite))
}

// InvitationAcceptanceSuite tests invitation acceptance endpoints
// POST /v1/invitations/accept and GET /v1/invitations/accept/google
type InvitationAcceptanceSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *InvitationAcceptanceSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *InvitationAcceptanceSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations")
	s.NoError(err)
}

func (s *InvitationAcceptanceSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations")
	s.NoError(err)
}

func (s *InvitationAcceptanceSuite) TestAcceptInvitation_Success() {
	ctx := context.Background()

	// Create workspace and invitation
	owner, workspace, _, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	invitation, token, err := s.helper.CreateInvitationWithToken(ctx, workspace.ID, "newuser@example.com", owner.ID, role.RoleUser)
	s.NoError(err)

	payload := map[string]interface{}{
		"firstName": "New",
		"lastName":  "User",
		"password":  "NewUserPassword123!",
	}

	resp, err := s.helper.Client.Post(fmt.Sprintf("/v1/invitations/accept?token=%s", token), payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))

	// Assert response structure
	s.Len(result, 2)
	s.Contains(result, "user")
	s.Contains(result, "token")
	s.NotEmpty(result["token"])

	userData := result["user"].(map[string]interface{})
	s.Equal("newuser@example.com", userData["email"])
	s.Equal("New", userData["firstName"])
	s.Equal("User", userData["lastName"])
	s.Equal(workspace.ID, userData["workspaceId"])
	s.Equal("user", userData["role"])
	s.Equal(true, userData["isEmailVerified"])

	// Verify invitation is marked as accepted
	updatedInvitation, err := s.helper.GetInvitation(ctx, invitation.ID)
	s.NoError(err)
	s.Equal(account.InvitationStatusAccepted, updatedInvitation.Status)
	s.NotNil(updatedInvitation.AcceptedAt)

	// Verify user can login with new credentials
	loginPayload := map[string]interface{}{
		"email":    "newuser@example.com",
		"password": "NewUserPassword123!",
	}
	loginResp, err := s.helper.Client.Post("/v1/auth/login", loginPayload)
	s.NoError(err)
	defer loginResp.Body.Close()
	s.Equal(http.StatusOK, loginResp.StatusCode)
}

func (s *InvitationAcceptanceSuite) TestAcceptInvitation_InvalidToken() {
	payload := map[string]interface{}{
		"firstName": "New",
		"lastName":  "User",
		"password":  "NewUserPassword123!",
	}

	resp, err := s.helper.Client.Post("/v1/invitations/accept?token=invalid-token", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result["detail"], "invalid or expired token")
}

func (s *InvitationAcceptanceSuite) TestAcceptInvitation_MissingToken() {
	payload := map[string]interface{}{
		"firstName": "New",
		"lastName":  "User",
		"password":  "NewUserPassword123!",
	}

	resp, err := s.helper.Client.Post("/v1/invitations/accept", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *InvitationAcceptanceSuite) TestAcceptInvitation_MissingPassword() {
	ctx := context.Background()
	owner, workspace, _, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	_, token, err := s.helper.CreateInvitationWithToken(ctx, workspace.ID, "newuser@example.com", owner.ID, role.RoleUser)
	s.NoError(err)

	payload := map[string]interface{}{
		"firstName": "New",
		"lastName":  "User",
	}

	resp, err := s.helper.Client.Post(fmt.Sprintf("/v1/invitations/accept?token=%s", token), payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *InvitationAcceptanceSuite) TestAcceptInvitation_ShortPassword() {
	ctx := context.Background()
	owner, workspace, _, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	_, token, err := s.helper.CreateInvitationWithToken(ctx, workspace.ID, "newuser@example.com", owner.ID, role.RoleUser)
	s.NoError(err)

	payload := map[string]interface{}{
		"firstName": "New",
		"lastName":  "User",
		"password":  "Short1!",
	}

	resp, err := s.helper.Client.Post(fmt.Sprintf("/v1/invitations/accept?token=%s", token), payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *InvitationAcceptanceSuite) TestAcceptInvitation_AlreadyAccepted() {
	ctx := context.Background()
	owner, workspace, _, err := s.helper.CreateTestUser(ctx, "owner@example.com", "Password123!", "Owner", "User", role.RoleAdmin)
	s.NoError(err)

	invitation, token, err := s.helper.CreateInvitationWithToken(ctx, workspace.ID, "newuser@example.com", owner.ID, role.RoleUser)
	s.NoError(err)

	// Mark invitation as accepted
	err = s.helper.SetInvitationStatus(ctx, invitation.ID, account.InvitationStatusAccepted)
	s.NoError(err)

	payload := map[string]interface{}{
		"firstName": "New",
		"lastName":  "User",
		"password":  "NewUserPassword123!",
	}

	resp, err := s.helper.Client.Post(fmt.Sprintf("/v1/invitations/accept?token=%s", token), payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusConflict, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result["detail"], "already been accepted")
}

func TestInvitationAcceptanceSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(InvitationAcceptanceSuite))
}

// UserManagementSuite tests user management endpoints
// PATCH /v1/workspaces/users/:userId/role and DELETE /v1/workspaces/users/:userId
type UserManagementSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *UserManagementSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *UserManagementSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations")
	s.NoError(err)
}

func (s *UserManagementSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations")
	s.NoError(err)
}

func (s *UserManagementSuite) TestUpdateUserRole_Success() {
	ctx := context.Background()

	// Create workspace with multiple users
	workspace, users, err := testutils.CreateWorkspaceWithUsers(ctx, testEnv.Database, "owner@example.com", "Password123!", []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Role      role.Role
	}{
		{"member@example.com", "Password123!", "Member", "User", role.RoleUser},
	})
	s.NoError(err)

	// Login as owner
	token, err := testutils.LoginAndGetToken(s.helper.Client, "owner@example.com", "Password123!")
	s.NoError(err)

	// Update member role to admin
	targetUser := users[1]
	payload := map[string]interface{}{
		"role": "admin",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/workspaces/users/"+targetUser.ID+"/role", payload, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal("admin", result["role"])
	s.Equal(targetUser.ID, result["id"])
	_ = workspace
}

func (s *UserManagementSuite) TestUpdateUserRole_CannotUpdateOwn() {
	ctx := context.Background()
	user, _, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)

	payload := map[string]interface{}{
		"role": "user",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/workspaces/users/"+user.ID+"/role", payload, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusForbidden, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result["detail"], "cannot update your own role")
}

func (s *UserManagementSuite) TestUpdateUserRole_CannotUpdateOwner() {
	ctx := context.Background()
	workspace, users, err := testutils.CreateWorkspaceWithUsers(ctx, testEnv.Database, "owner@example.com", "Password123!", []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Role      role.Role
	}{
		{"admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin},
	})
	s.NoError(err)

	// Login as admin (not owner)
	token, err := testutils.LoginAndGetToken(s.helper.Client, "admin@example.com", "Password123!")
	s.NoError(err)

	// Try to update owner's role
	owner := users[0]
	payload := map[string]interface{}{
		"role": "user",
	}

	resp, err := s.helper.Client.AuthenticatedRequest("PATCH", "/v1/workspaces/users/"+owner.ID+"/role", payload, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusForbidden, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result["detail"], "cannot update the workspace owner's role")
	_ = workspace
}

func (s *UserManagementSuite) TestRemoveUser_Success() {
	ctx := context.Background()
	workspace, users, err := testutils.CreateWorkspaceWithUsers(ctx, testEnv.Database, "owner@example.com", "Password123!", []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Role      role.Role
	}{
		{"member@example.com", "Password123!", "Member", "User", role.RoleUser},
	})
	s.NoError(err)

	// Login as owner
	token, err := testutils.LoginAndGetToken(s.helper.Client, "owner@example.com", "Password123!")
	s.NoError(err)

	// Remove member
	targetUser := users[1]
	resp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/workspaces/users/"+targetUser.ID, nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify user count decreased
	count, err := s.helper.CountWorkspaceUsers(ctx, workspace.ID)
	s.NoError(err)
	s.Equal(int64(1), count) // Only owner remains
}

func (s *UserManagementSuite) TestRemoveUser_CannotRemoveSelf() {
	ctx := context.Background()
	user, _, token, err := s.helper.CreateTestUser(ctx, "admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin)
	s.NoError(err)

	resp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/workspaces/users/"+user.ID, nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusForbidden, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result["detail"], "cannot remove yourself")
}

func (s *UserManagementSuite) TestRemoveUser_CannotRemoveOwner() {
	ctx := context.Background()
	_, users, err := testutils.CreateWorkspaceWithUsers(ctx, testEnv.Database, "owner@example.com", "Password123!", []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Role      role.Role
	}{
		{"admin@example.com", "Password123!", "Admin", "User", role.RoleAdmin},
	})
	s.NoError(err)

	// Login as admin
	token, err := testutils.LoginAndGetToken(s.helper.Client, "admin@example.com", "Password123!")
	s.NoError(err)

	// Try to remove owner
	owner := users[0]
	resp, err := s.helper.Client.AuthenticatedRequest("DELETE", "/v1/workspaces/users/"+owner.ID, nil, token)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusForbidden, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result["detail"], "cannot remove the workspace owner")
}

func TestUserManagementSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(UserManagementSuite))
}
