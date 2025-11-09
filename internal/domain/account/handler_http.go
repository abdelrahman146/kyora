package account

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	service *Service
}

func NewHttpHandler(service *Service) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

func (h *HttpHandler) InviteUserToWorkspace(c *gin.Context) {
	actor, err := ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	workspace, err := WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var input InviteUserInput
	if err := request.ValidBody(c, &input); err != nil {
		return
	}

	invitation, err := h.service.InviteUserToWorkspace(c.Request.Context(), actor, workspace.ID, input.Email, input.Role)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusCreated, invitation)
}

func (h *HttpHandler) GetWorkspaceInvitations(c *gin.Context) {
	workspace, err := WorkspaceFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	statusParam := c.Query("status")
	var status InvitationStatus
	if statusParam != "" {
		status = InvitationStatus(statusParam)
	}

	invitations, err := h.service.GetWorkspaceInvitations(c.Request.Context(), workspace.ID, status)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, invitations)
}

func (h *HttpHandler) RevokeInvitation(c *gin.Context) {
	invitationID := c.Param("invitationId")

	err := h.service.RevokeInvitation(c.Request.Context(), invitationID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}

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

func (h *HttpHandler) UpdateUserRole(c *gin.Context) {
	actor, err := ActorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	workspace, err := WorkspaceFromContext(c)
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

	updatedUser, err := h.service.UpdateUserRole(c.Request.Context(), actor, workspace, userID, input.Role)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessJSON(c, http.StatusOK, updatedUser)
}
