package account

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

// HttpHandler handles HTTP requests for account domain operations
type HttpHandler struct {
	service *Service
}

// NewHttpHandler creates a new HTTP handler for account operations
func NewHttpHandler(service *Service) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

// Authentication endpoints

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates a user with email and password
//
// @Summary      Login with email and password
// @Description  Authenticates a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body loginRequest true "Login credentials"
// @Success      200 {object} LoginResponse
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/auth/login [post]
func (h *HttpHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	loginResp, err := h.service.LoginWithEmailAndPasswordWithContext(c.Request.Context(), req.Email, req.Password, clientIP, userAgent)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, loginResp)
}

// GetGoogleAuthURL returns the Google OAuth URL for authentication
//
// @Summary      Get Google OAuth URL
// @Description  Returns the Google OAuth authorization URL for user authentication
// @Tags         auth
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      500 {object} problem.Problem
// @Router       /v1/auth/google/url [get]
func (h *HttpHandler) GetGoogleAuthURL(c *gin.Context) {
	url, state, err := h.service.GetGoogleAuthURL(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, gin.H{
		"url":   url,
		"state": state,
	})
}

type googleLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// LoginWithGoogle authenticates a user with Google OAuth
//
// @Summary      Login with Google OAuth
// @Description  Authenticates a user using Google OAuth code and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body googleLoginRequest true "Google OAuth code"
// @Success      200 {object} LoginResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/auth/google/login [post]
func (h *HttpHandler) LoginWithGoogle(c *gin.Context) {
	var req googleLoginRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	// Exchange Google code for user info
	googleUserInfo, err := h.service.ExchangeGoogleCodeAndFetchUser(c.Request.Context(), req.Code)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Try to find existing user
	user, err := h.service.GetUserByEmail(c.Request.Context(), googleUserInfo.Email)
	if err != nil {
		response.Error(c, problem.Unauthorized("no account found with this Google email").WithError(err))
		return
	}

	// Generate JWT token
	jwtToken, err := auth.NewJwtToken(user.ID, user.WorkspaceID)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	loginResp := &LoginResponse{
		User:  user,
		Token: jwtToken,
	}

	response.SuccessJSON(c, http.StatusOK, loginResp)
}

// Password reset endpoints

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgotPassword initiates the password reset process
//
// @Summary      Forgot password
// @Description  Sends a password reset email to the user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body forgotPasswordRequest true "User email"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/auth/forgot-password [post]
func (h *HttpHandler) ForgotPassword(c *gin.Context) {
	var req forgotPasswordRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	_, err := h.service.CreatePasswordResetToken(c.Request.Context(), req.Email)
	if err != nil {
		// Return success even if email not found to prevent email enumeration
		response.SuccessEmpty(c, http.StatusNoContent)
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

type resetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// ResetPassword resets the user's password using a reset token
//
// @Summary      Reset password
// @Description  Resets user password using a valid reset token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body resetPasswordRequest true "Reset password data"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/auth/reset-password [post]
func (h *HttpHandler) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

// Email verification endpoints

type requestEmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// RequestEmailVerification sends an email verification link
//
// @Summary      Request email verification
// @Description  Sends an email verification link to the user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body requestEmailVerificationRequest true "User email"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/auth/verify-email/request [post]
func (h *HttpHandler) RequestEmailVerification(c *gin.Context) {
	var req requestEmailVerificationRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	_, err := h.service.CreateVerifyEmailToken(c.Request.Context(), req.Email)
	if err != nil {
		// Return success even if email not found to prevent email enumeration
		response.SuccessEmpty(c, http.StatusNoContent)
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

type verifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyEmail verifies a user's email address
//
// @Summary      Verify email
// @Description  Verifies user's email address using a verification token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body verifyEmailRequest true "Verification token"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/auth/verify-email [post]
func (h *HttpHandler) VerifyEmail(c *gin.Context) {
	var req verifyEmailRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	err := h.service.VerifyEmail(c.Request.Context(), req.Token)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

// User profile endpoints

// GetCurrentUser returns the authenticated user's profile
//
// @Summary      Get current user
// @Description  Returns the profile of the currently authenticated user
// @Tags         users
// @Produce      json
// @Success      200 {object} User
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/users/me [get]
// @Security     BearerAuth
func (h *HttpHandler) GetCurrentUser(c *gin.Context) {
	actor, err := ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, actor)
}

// UpdateCurrentUser updates the authenticated user's profile
//
// @Summary      Update current user
// @Description  Updates the profile of the currently authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body UpdateUserInput true "User update data"
// @Success      200 {object} User
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/users/me [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateCurrentUser(c *gin.Context) {
	actor, err := ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input UpdateUserInput
	if err := request.ValidBody(c, &input); err != nil {
		return
	}

	// Update user fields
	if input.FirstName != nil {
		actor.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		actor.LastName = *input.LastName
	}

	// Persist changes
	updatedUser, err := h.service.GetUserByID(c.Request.Context(), actor.ID)
	if err != nil {
		response.Error(c, err)
		return
	}

	if input.FirstName != nil {
		updatedUser.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		updatedUser.LastName = *input.LastName
	}

	if err := h.service.storage.user.UpdateOne(c.Request.Context(), updatedUser); err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, updatedUser)
}

// Workspace endpoints

// GetCurrentWorkspace returns the authenticated user's workspace
//
// @Summary      Get current workspace
// @Description  Returns the workspace of the currently authenticated user
// @Tags         workspaces
// @Produce      json
// @Success      200 {object} Workspace
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/workspaces/me [get]
// @Security     BearerAuth
func (h *HttpHandler) GetCurrentWorkspace(c *gin.Context) {
	actor, err := ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	workspace, err := h.service.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, workspace)
}

// GetWorkspaceUsers returns all users in the workspace
//
// @Summary      Get workspace users
// @Description  Returns all users that belong to the workspace
// @Tags         workspaces
// @Produce      json
// @Success      200 {array} User
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/workspaces/users [get]
// @Security     BearerAuth
func (h *HttpHandler) GetWorkspaceUsers(c *gin.Context) {
	actor, err := ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	users, err := h.service.storage.user.FindMany(
		c.Request.Context(),
		h.service.storage.user.ScopeWorkspaceID(actor.WorkspaceID),
		h.service.storage.user.WithOrderBy([]string{UserSchema.CreatedAt.Column() + " ASC"}),
	)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessJSON(c, http.StatusOK, users)
}

// GetWorkspaceUser returns a specific user in the workspace
//
// @Summary      Get workspace user
// @Description  Returns a specific user by ID from the workspace
// @Tags         workspaces
// @Produce      json
// @Param        userId path string true "User ID"
// @Success      200 {object} User
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/workspaces/users/{userId} [get]
// @Security     BearerAuth
func (h *HttpHandler) GetWorkspaceUser(c *gin.Context) {
	actor, err := ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	userID := c.Param("userId")
	if userID == "" {
		response.Error(c, problem.BadRequest("userId is required"))
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Verify user belongs to same workspace
	if user.WorkspaceID != actor.WorkspaceID {
		response.Error(c, ErrUserNotInWorkspace(nil))
		return
	}

	response.SuccessJSON(c, http.StatusOK, user)
}

// Invitation endpoints

// InviteUserToWorkspace invites a user to join the workspace
//
// @Summary      Invite user to workspace
// @Description  Sends an invitation email to a user to join the workspace
// @Tags         invitations
// @Accept       json
// @Produce      json
// @Param        request body InviteUserInput true "Invitation data"
// @Success      201 {object} UserInvitation
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/invitations [post]
// @Security     BearerAuth
func (h *HttpHandler) InviteUserToWorkspace(c *gin.Context) {
	actor, err := ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input InviteUserInput
	if err := request.ValidBody(c, &input); err != nil {
		return
	}

	invitation, err := h.service.InviteUserToWorkspace(c.Request.Context(), actor, actor.WorkspaceID, input.Email, input.Role)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusCreated, invitation)
}

// GetWorkspaceInvitations returns all invitations for the workspace
//
// @Summary      Get workspace invitations
// @Description  Returns all invitations for the workspace, optionally filtered by status
// @Tags         invitations
// @Produce      json
// @Param        status query string false "Filter by invitation status (pending, accepted, expired, revoked)"
// @Success      200 {array} UserInvitation
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/invitations [get]
// @Security     BearerAuth
func (h *HttpHandler) GetWorkspaceInvitations(c *gin.Context) {
	actor, err := ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	statusParam := c.Query("status")
	var status InvitationStatus
	if statusParam != "" {
		status = InvitationStatus(statusParam)
	}

	invitations, err := h.service.GetWorkspaceInvitations(c.Request.Context(), actor.WorkspaceID, status)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, invitations)
}

// RevokeInvitation revokes a pending invitation
//
// @Summary      Revoke invitation
// @Description  Revokes a pending workspace invitation
// @Tags         invitations
// @Produce      json
// @Param        invitationId path string true "Invitation ID"
// @Success      204
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/invitations/{invitationId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) RevokeInvitation(c *gin.Context) {
	invitationID := c.Param("invitationId")
	if invitationID == "" {
		response.Error(c, problem.BadRequest("invitationId is required"))
		return
	}

	err := h.service.RevokeInvitation(c.Request.Context(), invitationID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

// AcceptInvitation accepts a workspace invitation (public endpoint)
//
// @Summary      Accept invitation
// @Description  Accepts a workspace invitation and creates a new user account
// @Tags         invitations
// @Accept       json
// @Produce      json
// @Param        token query string true "Invitation token"
// @Param        request body CreateUserInput true "User account data"
// @Success      200 {object} LoginResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/invitations/accept [post]
func (h *HttpHandler) AcceptInvitation(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.Error(c, problem.BadRequest("token is required"))
		return
	}

	var input CreateUserInput
	if err := request.ValidBody(c, &input); err != nil {
		return
	}

	user, workspace, err := h.service.AcceptInvitation(c.Request.Context(), token, input.FirstName, input.LastName, input.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Load workspace relationship
	user.Workspace = workspace

	// Generate JWT token for auto-login
	jwtToken, err := auth.NewJwtToken(user.ID, user.WorkspaceID)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	loginResp := &LoginResponse{
		User:  user,
		Token: jwtToken,
	}

	response.SuccessJSON(c, http.StatusOK, loginResp)
}

// AcceptInvitationWithGoogle accepts a workspace invitation using Google OAuth (public endpoint)
//
// @Summary      Accept invitation with Google
// @Description  Accepts a workspace invitation and creates a new user account using Google OAuth
// @Tags         invitations
// @Accept       json
// @Produce      json
// @Param        token query string true "Invitation token"
// @Param        code query string true "Google OAuth code"
// @Success      200 {object} LoginResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/invitations/accept/google [get]
func (h *HttpHandler) AcceptInvitationWithGoogle(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.Error(c, problem.BadRequest("token is required"))
		return
	}

	code := c.Query("code")
	if code == "" {
		response.Error(c, problem.BadRequest("code is required"))
		return
	}

	// Exchange Google code for user info
	googleUserInfo, err := h.service.ExchangeGoogleCodeAndFetchUser(c.Request.Context(), code)
	if err != nil {
		response.Error(c, err)
		return
	}

	user, workspace, err := h.service.AcceptInvitationWithGoogleAuth(c.Request.Context(), token, googleUserInfo)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Load workspace relationship
	user.Workspace = workspace

	// Generate JWT token for auto-login
	jwtToken, err := auth.NewJwtToken(user.ID, user.WorkspaceID)
	if err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	loginResp := &LoginResponse{
		User:  user,
		Token: jwtToken,
	}

	response.SuccessJSON(c, http.StatusOK, loginResp)
}

// User management endpoints

// UpdateUserRole updates a user's role within the workspace
//
// @Summary      Update user role
// @Description  Updates the role of a user within the workspace
// @Tags         workspaces
// @Accept       json
// @Produce      json
// @Param        userId path string true "User ID"
// @Param        request body UpdateUserRoleInput true "Role update data"
// @Success      200 {object} User
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/workspaces/users/{userId}/role [patch]
// @Security     BearerAuth
func (h *HttpHandler) UpdateUserRole(c *gin.Context) {
	actor, err := ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	userID := c.Param("userId")
	if userID == "" {
		response.Error(c, problem.BadRequest("userId is required"))
		return
	}

	var input UpdateUserRoleInput
	if err := request.ValidBody(c, &input); err != nil {
		return
	}

	// Get workspace from actor
	workspace, err := h.service.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
	if err != nil {
		response.Error(c, err)
		return
	}

	updatedUser, err := h.service.UpdateUserRole(c.Request.Context(), actor, workspace, userID, input.Role)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, updatedUser)
}

// RemoveUserFromWorkspace removes a user from the workspace
//
// @Summary      Remove user from workspace
// @Description  Removes a user from the workspace (soft delete)
// @Tags         workspaces
// @Produce      json
// @Param        userId path string true "User ID"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      403 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/workspaces/users/{userId} [delete]
// @Security     BearerAuth
func (h *HttpHandler) RemoveUserFromWorkspace(c *gin.Context) {
	actor, err := ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	userID := c.Param("userId")
	if userID == "" {
		response.Error(c, problem.BadRequest("userId is required"))
		return
	}

	// Prevent user from removing themselves
	if actor.ID == userID {
		response.Error(c, problem.Forbidden("you cannot remove yourself from the workspace"))
		return
	}

	// Get target user
	targetUser, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Verify target user is in the same workspace
	if targetUser.WorkspaceID != actor.WorkspaceID {
		response.Error(c, ErrUserNotInWorkspace(nil))
		return
	}

	// Get workspace
	workspace, err := h.service.GetWorkspaceByID(c.Request.Context(), actor.WorkspaceID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Prevent removing the workspace owner
	if targetUser.ID == workspace.OwnerID {
		response.Error(c, problem.Forbidden("you cannot remove the workspace owner"))
		return
	}

	// Soft delete the user
	if err := h.service.storage.user.DeleteOne(c.Request.Context(), targetUser); err != nil {
		response.Error(c, problem.InternalError().WithError(err))
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}
