package e2e_test

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"gorm.io/gorm"
)

// AccountTestHelper provides reusable helpers for account tests
type AccountTestHelper struct {
	db      *database.Database
	Client  *testutils.HTTPClient
	Storage *account.Storage
}

// NewAccountTestHelper creates a new account test helper
func NewAccountTestHelper(db *database.Database, cacheAddr, baseURL string) *AccountTestHelper {
	cacheClient := cache.NewConnection([]string{cacheAddr})
	storage := account.NewStorage(db, cacheClient)

	return &AccountTestHelper{
		db:      db,
		Client:  testutils.NewHTTPClient(baseURL),
		Storage: storage,
	}
}

// CreateTestUser creates a user and workspace using storage layer
func (h *AccountTestHelper) CreateTestUser(ctx context.Context, email, password, firstName, lastName string, userRole role.Role) (*account.User, *account.Workspace, string, error) {
	return testutils.CreateAuthenticatedUser(ctx, h.db, email, password, firstName, lastName, userRole)
}

// CreateInvitation creates a user invitation record using repository
func (h *AccountTestHelper) CreateInvitation(ctx context.Context, workspaceID, email string, inviterID string, userRole role.Role, status account.InvitationStatus) (*account.UserInvitation, error) {
	invitationRepo := database.NewRepository[account.UserInvitation](h.db)
	invitation := &account.UserInvitation{
		WorkspaceID: workspaceID,
		Email:       email,
		Role:        userRole,
		InviterID:   inviterID,
		Status:      status,
	}

	if err := invitationRepo.CreateOne(ctx, invitation); err != nil {
		return nil, err
	}

	return invitation, nil
}

// CreateInvitationWithToken creates an invitation and returns the token
func (h *AccountTestHelper) CreateInvitationWithToken(ctx context.Context, workspaceID, email string, inviterID string, userRole role.Role) (*account.UserInvitation, string, error) {
	invitation, err := h.CreateInvitation(ctx, workspaceID, email, inviterID, userRole, account.InvitationStatusPending)
	if err != nil {
		return nil, "", err
	}

	token, err := testutils.CreateInvitationToken(ctx, h.Storage, invitation, inviterID)
	if err != nil {
		return nil, "", err
	}

	return invitation, token, nil
}

// CreateInvitationToken generates a token for an existing invitation
func (h *AccountTestHelper) CreateInvitationToken(ctx context.Context, invitationID string) (string, error) {
	invitation, err := h.GetInvitation(ctx, invitationID)
	if err != nil {
		return "", err
	}

	return testutils.CreateInvitationToken(ctx, h.Storage, invitation, invitation.InviterID)
}

// SetInvitationStatus updates invitation status using repository
func (h *AccountTestHelper) SetInvitationStatus(ctx context.Context, invitationID string, status account.InvitationStatus) error {
	invitationRepo := database.NewRepository[account.UserInvitation](h.db)
	invitation, err := invitationRepo.FindByID(ctx, invitationID)
	if err != nil {
		return err
	}

	invitation.Status = status
	if status == account.InvitationStatusAccepted {
		now := time.Now()
		invitation.AcceptedAt = &gorm.DeletedAt{Time: now, Valid: true}
	}

	return invitationRepo.UpdateOne(ctx, invitation)
}

// GetUser retrieves a user by ID using repository
func (h *AccountTestHelper) GetUser(ctx context.Context, userID string) (*account.User, error) {
	userRepo := database.NewRepository[account.User](h.db)
	return userRepo.FindByID(ctx, userID)
}

// GetWorkspace retrieves a workspace by ID using repository
func (h *AccountTestHelper) GetWorkspace(ctx context.Context, workspaceID string) (*account.Workspace, error) {
	workspaceRepo := database.NewRepository[account.Workspace](h.db)
	return workspaceRepo.FindByID(ctx, workspaceID)
}

// GetInvitation retrieves an invitation by ID using repository
func (h *AccountTestHelper) GetInvitation(ctx context.Context, invitationID string) (*account.UserInvitation, error) {
	invitationRepo := database.NewRepository[account.UserInvitation](h.db)
	return invitationRepo.FindByID(ctx, invitationID)
}

// CreatePasswordResetToken creates a password reset token
func (h *AccountTestHelper) CreatePasswordResetToken(ctx context.Context, user *account.User) (string, error) {
	return testutils.CreatePasswordResetToken(ctx, h.Storage, user)
}

// CreateEmailVerificationToken creates an email verification token
func (h *AccountTestHelper) CreateEmailVerificationToken(ctx context.Context, user *account.User) (string, error) {
	return testutils.CreateEmailVerificationToken(ctx, h.Storage, user)
}

// MarkEmailUnverified marks a user's email as unverified
func (h *AccountTestHelper) MarkEmailUnverified(ctx context.Context, userID string) error {
	user, err := h.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	userRepo := database.NewRepository[account.User](h.db)
	user.IsEmailVerified = false
	return userRepo.UpdateOne(ctx, user)
}

// UpdateUserPassword updates a user's password
func (h *AccountTestHelper) UpdateUserPassword(ctx context.Context, userID, hashedPassword string) error {
	user, err := h.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	userRepo := database.NewRepository[account.User](h.db)
	user.Password = hashedPassword
	return userRepo.UpdateOne(ctx, user)
}

// CountWorkspaceUsers counts users in a workspace
func (h *AccountTestHelper) CountWorkspaceUsers(ctx context.Context, workspaceID string) (int64, error) {
	userRepo := database.NewRepository[account.User](h.db)
	return userRepo.Count(ctx, userRepo.ScopeWorkspaceID(workspaceID))
}

// ListWorkspaceInvitations lists all invitations for a workspace with optional status filter
func (h *AccountTestHelper) ListWorkspaceInvitations(ctx context.Context, workspaceID string, status account.InvitationStatus) ([]*account.UserInvitation, error) {
	invitationRepo := database.NewRepository[account.UserInvitation](h.db)
	scopes := []func(db *gorm.DB) *gorm.DB{
		invitationRepo.ScopeWorkspaceID(workspaceID),
	}
	if status != "" {
		scopes = append(scopes, invitationRepo.ScopeEquals(account.UserInvitationSchema.Status, string(status)))
	}
	return invitationRepo.FindMany(ctx, scopes...)
}

// CreateTestSubscription creates a test subscription for a workspace
func (h *AccountTestHelper) CreateTestSubscription(ctx context.Context, workspaceID string) error {
	return testutils.CreateTestSubscription(ctx, h.db, workspaceID)
}
