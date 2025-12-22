package account

import (
	"context"
	"fmt"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/email"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"gorm.io/gorm"
)

type Service struct {
	bus             *bus.Bus
	storage         *Storage
	atomicProcessor atomic.AtomicProcessor
	Notification    *Notification
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, bus *bus.Bus, emailClient email.Client) *Service {
	emailInfo := email.NewEmail()
	notification := NewNotification(emailClient, emailInfo)
	return &Service{
		storage:         storage,
		atomicProcessor: atomicProcessor,
		bus:             bus,
		Notification:    notification,
	}
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	return s.storage.user.FindByID(ctx, id, s.storage.user.WithPreload(WorkspaceStruct))
}

func (s *Service) GetWorkspaceByID(ctx context.Context, id string) (*Workspace, error) {
	return s.storage.workspace.FindByID(ctx, id, s.storage.workspace.WithPreload("Users"))
}

// SetWorkspaceStripeCustomer sets the Stripe customer ID for a workspace
func (s *Service) SetWorkspaceStripeCustomer(ctx context.Context, workspaceID, customerID string) error {
	ws, err := s.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return err
	}
	ws.StripeCustomerID.String = customerID
	ws.StripeCustomerID.Valid = true
	return s.storage.workspace.UpdateOne(ctx, ws)
}

// SetWorkspaceDefaultPaymentMethod sets the default Stripe payment method ID for a workspace
func (s *Service) SetWorkspaceDefaultPaymentMethod(ctx context.Context, workspaceID, pmID string) error {
	ws, err := s.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return err
	}
	ws.StripePaymentMethodID.String = pmID
	ws.StripePaymentMethodID.Valid = true
	return s.storage.workspace.UpdateOne(ctx, ws)
}

// CountWorkspaceUsers returns number of users in the workspace
func (s *Service) CountWorkspaceUsers(ctx context.Context, workspaceID string) (int64, error) {
	return s.storage.user.Count(ctx, s.storage.user.ScopeWorkspaceID(workspaceID))
}

// CountWorkspaceUsersForPlanLimit is a wrapper that matches the billing EnforcePlanLimitFunc signature
// It counts users in the given workspace (the id parameter is the workspaceID)
func (s *Service) CountWorkspaceUsersForPlanLimit(ctx context.Context, actor *User, id string) (int64, error) {
	return s.CountWorkspaceUsers(ctx, id)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.storage.user.FindOne(ctx, s.storage.user.ScopeEquals(UserSchema.Email, email), s.storage.user.WithPreload(WorkspaceStruct))
}

type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

func (s *Service) LoginWithEmailAndPassword(ctx context.Context, email, password string) (*LoginResponse, error) {
	return s.LoginWithEmailAndPasswordWithContext(ctx, email, password, "", "")
}

// LoginWithEmailAndPasswordWithContext includes client context for security notifications
func (s *Service) LoginWithEmailAndPasswordWithContext(ctx context.Context, email, password, clientIP, userAgent string) (*LoginResponse, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials(err)
	}
	if hash.ValidatePassword(password, user.Password) {
		token, err := auth.NewJwtToken(user.ID, user.WorkspaceID)
		if err != nil {
			return nil, problem.InternalError().WithError(err)
		}

		// Send login notification email asynchronously
		go func() {
			if err := s.Notification.SendLoginNotificationEmail(context.Background(), user, clientIP, userAgent); err != nil {
				// Log error but don't fail the login process
				// Login notifications are a security feature but shouldn't block user access
			}
		}()

		return &LoginResponse{
			User:  user,
			Token: token,
		}, nil
	}
	return nil, ErrInvalidCredentials(nil)
}

func (s *Service) CreateVerifyEmailToken(ctx context.Context, email string) (string, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials(err)
	}
	token, expAt, err := s.storage.CreateVerifyEmailToken(&VerifyEmailPayload{
		UserID:      user.ID,
		WorkspaceID: user.WorkspaceID,
		Email:       user.Email,
	})
	if err != nil {
		return "", problem.InternalError().WithError(err)
	}

	err = s.Notification.SendEmailVerificationEmail(ctx, user, token, expAt)
	if err != nil {
		// Log error but don't fail the token creation
	}
	return token, nil
}

func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	payload, err := s.storage.GetVerifyEmailToken(token)
	if err != nil {
		return ErrInvalidOrExpiredToken(err)
	}
	user, err := s.GetUserByID(ctx, payload.UserID)
	if err != nil {
		return ErrInvalidOrExpiredToken(err)
	}
	user.IsEmailVerified = true
	err = s.storage.user.UpdateOne(ctx, user)
	if err != nil {
		return problem.InternalError().WithError(err)
	}
	err = s.storage.ConsumeVerifyEmailToken(token)
	if err != nil {
		return problem.InternalError().WithError(err)
	}
	return nil
}

func (s *Service) CreatePasswordResetToken(ctx context.Context, email string) (string, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials(err)
	}
	token, expAt, err := s.storage.CreatePasswordResetToken(&PasswordResetPayload{
		UserID:      user.ID,
		WorkspaceID: user.WorkspaceID,
		Email:       user.Email,
	})
	if err != nil {
		return "", problem.InternalError().WithError(err)
	}

	// Try to send email using the email integration
	err = s.Notification.SendForgotPasswordEmail(ctx, user, token, expAt)
	if err != nil {
		// Log error but don't fail the token creation
		logger.FromContext(ctx).Error("Failed to send forgot password email", "error", err)
	}
	return token, nil
}

func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	payload, err := s.storage.GetPasswordResetToken(token)
	if err != nil {
		return ErrInvalidOrExpiredToken(err)
	}
	user, err := s.GetUserByID(ctx, payload.UserID)
	if err != nil {
		return ErrInvalidOrExpiredToken(err)
	}
	hashedPassword, err := hash.Password(newPassword)
	if err != nil {
		return problem.InternalError().WithError(err)
	}
	user.Password = hashedPassword
	err = s.storage.user.UpdateOne(ctx, user)
	if err != nil {
		return problem.InternalError().WithError(err)
	}
	err = s.storage.ConsumePasswordResetToken(token)
	if err != nil {
		return problem.InternalError().WithError(err)
	}

	// Send password reset confirmation email
	err = s.Notification.SendPasswordResetConfirmationEmail(ctx, user, "")
	if err != nil {
		// Log error but don't fail the password reset operation
		// The password has already been changed successfully
	}
	return nil
}

func (s *Service) GetGoogleAuthURL(ctx context.Context) (url string, state string, err error) {
	state, err = id.RandomString(24)
	if err != nil {
		return "", "", problem.InternalError().WithError(err)
	}
	url, err = auth.GoogleGetAuthURL(ctx, state)
	if err != nil {
		return "", "", problem.InternalError().WithError(err)
	}
	return url, state, nil
}

func (s *Service) ExchangeGoogleCodeAndFetchUser(ctx context.Context, code string) (*auth.GoogleUserInfo, error) {
	info, err := auth.GoogleExchangeAndFetchUser(ctx, code)
	if err != nil {
		return nil, problem.InternalError().WithError(err)
	}
	return info, nil
}

// BootstrapWorkspaceAndOwner creates a new workspace and an owner user atomically.
// It avoids exposing storage details to callers that need to initialize a tenant.
func (s *Service) BootstrapWorkspaceAndOwner(ctx context.Context, firstName, lastName, email, passwordHash string, emailVerified bool, stripeCustomerID string) (*User, *Workspace, error) {
	var createdUser *User
	var createdWs *Workspace
	err := s.atomicProcessor.Exec(ctx, func(txCtx context.Context) error {
		// Guard against existing user
		if u, err := s.GetUserByEmail(txCtx, email); err == nil && u != nil {
			return ErrInvalidCredentials(nil)
		}
		ws := &Workspace{}
		if err := s.storage.workspace.CreateOne(txCtx, ws); err != nil {
			return err
		}
		user := &User{
			WorkspaceID:     ws.ID,
			Role:            role.RoleAdmin,
			FirstName:       firstName,
			LastName:        lastName,
			Email:           email,
			Password:        passwordHash,
			IsEmailVerified: emailVerified,
		}
		if err := s.storage.user.CreateOne(txCtx, user); err != nil {
			return err
		}
		ws.OwnerID = user.ID
		if stripeCustomerID != "" {
			ws.StripeCustomerID.String = stripeCustomerID
			ws.StripeCustomerID.Valid = true
		}
		if err := s.storage.workspace.UpdateOne(txCtx, ws); err != nil {
			return err
		}
		createdUser = user
		createdWs = ws
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return createdUser, createdWs, nil
}

// InviteUserToWorkspace creates a workspace invitation and sends an invitation email
func (s *Service) InviteUserToWorkspace(ctx context.Context, actor *User, workspaceID, email string, role role.Role) (*UserInvitation, error) {
	// Check if user already exists in the workspace
	existingUser, err := s.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		if existingUser.WorkspaceID == workspaceID {
			return nil, ErrUserAlreadyExists(nil)
		}
	}

	// Check if there's already a pending invitation
	existingInvitation, err := s.storage.invitation.FindOne(
		ctx,
		s.storage.invitation.ScopeWorkspaceID(workspaceID),
		s.storage.invitation.ScopeEquals(UserInvitationSchema.Email, email),
		s.storage.invitation.ScopeEquals(UserInvitationSchema.Status, string(InvitationStatusPending)),
	)
	if err == nil && existingInvitation != nil {
		return nil, ErrInvitationAlreadyExists(nil)
	}

	// Get workspace details for email
	workspace, err := s.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	// Create the invitation record
	invitation := &UserInvitation{
		WorkspaceID: workspaceID,
		Email:       email,
		Role:        role,
		InviterID:   actor.ID,
		Status:      InvitationStatusPending,
	}

	if err := s.storage.invitation.CreateOne(ctx, invitation); err != nil {
		return nil, problem.InternalError().WithError(err)
	}

	// Create invitation token
	token, expAt, err := s.storage.CreateWorkspaceInvitationToken(&WorkspaceInvitationPayload{
		InvitationID: invitation.ID,
		WorkspaceID:  workspaceID,
		Email:        email,
		Role:         string(role),
		InviterID:    actor.ID,
	})
	if err != nil {
		return nil, problem.InternalError().WithError(err)
	}

	// Send invitation email
	inviterName := fmt.Sprintf("%s %s", actor.FirstName, actor.LastName)
	if inviterName == " " {
		inviterName = actor.Email
	}
	workspaceName := fmt.Sprintf("Workspace %s", workspace.ID) // You may want to add a Name field to Workspace

	err = s.Notification.SendWorkspaceInvitationEmail(ctx, email, workspaceName, inviterName, actor.Email, string(role), token, expAt)
	if err != nil {
		// Log error but don't fail
		logger.FromContext(ctx).Error("Failed to send workspace invitation email, falling back to bus", "error", err)
	}

	return invitation, nil
}

// AcceptInvitation processes an invitation acceptance and creates a user account
func (s *Service) AcceptInvitation(ctx context.Context, token string, firstName, lastName, password string) (*User, *Workspace, error) {
	// Get invitation payload from token
	payload, err := s.storage.GetWorkspaceInvitationToken(token)
	if err != nil {
		return nil, nil, ErrInvalidOrExpiredToken(err)
	}

	// Find the invitation
	invitation, err := s.storage.invitation.FindByID(ctx, payload.InvitationID)
	if err != nil {
		return nil, nil, ErrInvitationNotFound(err)
	}

	// Check invitation status
	if invitation.Status == InvitationStatusAccepted {
		return nil, nil, ErrInvitationAlreadyAccepted(nil)
	}
	if invitation.Status == InvitationStatusExpired || invitation.Status == InvitationStatusRevoked {
		return nil, nil, ErrInvitationExpired(nil)
	}

	// Check if user already exists
	existingUser, err := s.GetUserByEmail(ctx, invitation.Email)
	if err == nil && existingUser != nil {
		return nil, nil, ErrUserAlreadyExists(nil)
	}

	// Hash password
	hashedPassword, err := hash.Password(password)
	if err != nil {
		return nil, nil, problem.InternalError().WithError(err)
	}

	var createdUser *User
	var workspace *Workspace

	// Create user and update invitation atomically
	err = s.atomicProcessor.Exec(ctx, func(txCtx context.Context) error {
		// Create user
		user := &User{
			WorkspaceID:     invitation.WorkspaceID,
			Role:            invitation.Role,
			FirstName:       firstName,
			LastName:        lastName,
			Email:           invitation.Email,
			Password:        hashedPassword,
			IsEmailVerified: true, // Auto-verify email for invited users
		}
		if err := s.storage.user.CreateOne(txCtx, user); err != nil {
			return err
		}
		createdUser = user

		// Update invitation status
		now := gorm.DeletedAt{Time: time.Now(), Valid: true}
		invitation.Status = InvitationStatusAccepted
		invitation.AcceptedAt = &now
		if err := s.storage.invitation.UpdateOne(txCtx, invitation); err != nil {
			return err
		}

		// Get workspace
		ws, err := s.GetWorkspaceByID(txCtx, invitation.WorkspaceID)
		if err != nil {
			return err
		}
		workspace = ws

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	// Consume the token
	if err := s.storage.ConsumeWorkspaceInvitationToken(token); err != nil {
		// Log but don't fail
		logger.FromContext(ctx).Error("Failed to consume workspace invitation token", "error", err)
	}

	return createdUser, workspace, nil
}

// AcceptInvitationWithGoogleAuth processes invitation acceptance for Google OAuth users
func (s *Service) AcceptInvitationWithGoogleAuth(ctx context.Context, token string, googleUserInfo *auth.GoogleUserInfo) (*User, *Workspace, error) {
	// Get invitation payload from token
	payload, err := s.storage.GetWorkspaceInvitationToken(token)
	if err != nil {
		return nil, nil, ErrInvalidOrExpiredToken(err)
	}

	// Verify the Google email matches the invitation email
	if googleUserInfo.Email != payload.Email {
		return nil, nil, problem.Forbidden("Google account email does not match the invitation email")
	}

	// Find the invitation
	invitation, err := s.storage.invitation.FindByID(ctx, payload.InvitationID)
	if err != nil {
		return nil, nil, ErrInvitationNotFound(err)
	}

	// Check invitation status
	if invitation.Status == InvitationStatusAccepted {
		return nil, nil, ErrInvitationAlreadyAccepted(nil)
	}
	if invitation.Status == InvitationStatusExpired || invitation.Status == InvitationStatusRevoked {
		return nil, nil, ErrInvitationExpired(nil)
	}

	// Check if user already exists
	existingUser, err := s.GetUserByEmail(ctx, invitation.Email)
	if err == nil && existingUser != nil {
		return nil, nil, ErrUserAlreadyExists(nil)
	}

	var createdUser *User
	var workspace *Workspace

	// Create user and update invitation atomically
	err = s.atomicProcessor.Exec(ctx, func(txCtx context.Context) error {
		// Create user without password (Google OAuth)
		user := &User{
			WorkspaceID:     invitation.WorkspaceID,
			Role:            invitation.Role,
			FirstName:       googleUserInfo.GivenName,
			LastName:        googleUserInfo.FamilyName,
			Email:           googleUserInfo.Email,
			Password:        "", // No password for Google OAuth users
			IsEmailVerified: googleUserInfo.Verified,
		}
		if err := s.storage.user.CreateOne(txCtx, user); err != nil {
			return err
		}
		createdUser = user

		// Update invitation status
		now := gorm.DeletedAt{Time: time.Now(), Valid: true}
		invitation.Status = InvitationStatusAccepted
		invitation.AcceptedAt = &now
		if err := s.storage.invitation.UpdateOne(txCtx, invitation); err != nil {
			return err
		}

		// Get workspace
		ws, err := s.GetWorkspaceByID(txCtx, invitation.WorkspaceID)
		if err != nil {
			return err
		}
		workspace = ws

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	// Consume the token
	if err := s.storage.ConsumeWorkspaceInvitationToken(token); err != nil {
		// Log but don't fail
	}

	return createdUser, workspace, nil
}

// UpdateUserRole updates a user's role within a workspace
func (s *Service) UpdateUserRole(ctx context.Context, actor *User, workspace *Workspace, targetUserID string, newRole role.Role) (*User, error) {
	// Prevent user from updating their own role
	if actor.ID == targetUserID {
		return nil, ErrCannotUpdateOwnRole(nil)
	}

	// Get target user
	targetUser, err := s.GetUserByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	// Verify target user is in the same workspace
	if targetUser.WorkspaceID != workspace.ID {
		return nil, ErrUserNotInWorkspace(nil)
	}

	// Prevent updating the workspace owner's role
	if targetUser.ID == workspace.OwnerID {
		return nil, ErrCannotUpdateOwnerRole(nil)
	}

	// Update role
	targetUser.Role = newRole
	if err := s.storage.user.UpdateOne(ctx, targetUser); err != nil {
		return nil, problem.InternalError().WithError(err)
	}

	return targetUser, nil
}

// GetWorkspaceInvitations returns all invitations for a workspace
func (s *Service) GetWorkspaceInvitations(ctx context.Context, workspaceID string, status InvitationStatus) ([]*UserInvitation, error) {
	var scopes []func(db *gorm.DB) *gorm.DB
	scopes = append(scopes, s.storage.invitation.ScopeWorkspaceID(workspaceID))
	scopes = append(scopes, s.storage.invitation.WithPreload("Inviter"))

	if status != "" {
		scopes = append(scopes, s.storage.invitation.ScopeEquals(UserInvitationSchema.Status, string(status)))
	}

	return s.storage.invitation.FindMany(ctx, scopes...)
}

// RevokeInvitation revokes a pending invitation
func (s *Service) RevokeInvitation(ctx context.Context, invitationID string) error {
	invitation, err := s.storage.invitation.FindByID(ctx, invitationID)
	if err != nil {
		return ErrInvitationNotFound(err)
	}

	if invitation.Status != InvitationStatusPending {
		return problem.Conflict("only pending invitations can be revoked")
	}

	invitation.Status = InvitationStatusRevoked
	return s.storage.invitation.UpdateOne(ctx, invitation)
}
