package account

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/config"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/email"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/abdelrahman146/kyora/internal/platform/utils/throttle"
	"github.com/spf13/viper"
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

// GetWorkspaceUserByID returns a user only if they belong to the given workspace.
// This prevents cross-workspace user ID probing (BOLA) and keeps errors consistent.
func (s *Service) GetWorkspaceUserByID(ctx context.Context, workspaceID, userID string) (*User, error) {
	user, err := s.storage.user.FindOne(
		ctx,
		s.storage.user.ScopeWorkspaceID(workspaceID),
		s.storage.user.ScopeID(userID),
		s.storage.user.WithPreload(WorkspaceStruct),
	)
	if err != nil {
		if database.IsRecordNotFound(err) {
			return nil, ErrUserNotInWorkspace(err)
		}
		return nil, problem.InternalError().WithError(err)
	}
	return user, nil
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

func (s *Service) issueTokensForUser(ctx context.Context, user *User, clientIP, userAgent string) (*RefreshResponse, error) {
	// Ensure we always have a non-zero auth version moving forward.
	if user.AuthVersion <= 0 {
		user.AuthVersion = 1
		if err := s.storage.user.UpdateOne(ctx, user); err != nil {
			return nil, problem.InternalError().WithError(err)
		}
	}

	accessToken, err := auth.NewJwtToken(user.ID, user.WorkspaceID, user.AuthVersion)
	if err != nil {
		return nil, problem.InternalError().WithError(err)
	}

	rawRefresh, err := auth.NewRefreshToken()
	if err != nil {
		return nil, problem.InternalError().WithError(err)
	}
	hash := auth.HashRefreshToken(rawRefresh)

	ttl := viper.GetInt(config.RefreshTokenExpirySeconds)
	if ttl <= 0 {
		ttl = 30 * 24 * 60 * 60
	}
	expAt := time.Now().UTC().Add(time.Duration(ttl) * time.Second)

	sess := &Session{
		UserID:      user.ID,
		WorkspaceID: user.WorkspaceID,
		TokenHash:   hash,
		ExpiresAt:   expAt,
		CreatedIP:   strings.TrimSpace(clientIP),
		UserAgent:   strings.TrimSpace(userAgent),
	}
	if sess.CreatedIP == "" {
		sess.CreatedIP = "unknown"
	}
	if err := s.storage.CreateSession(ctx, sess); err != nil {
		return nil, problem.InternalError().WithError(err)
	}

	return &RefreshResponse{Token: accessToken, RefreshToken: rawRefresh}, nil
}

// IssueTokensForUserWithContext issues a new access JWT and a new refresh token.
// The refresh token is stored hashed in the database and returned raw only once.
func (s *Service) IssueTokensForUserWithContext(ctx context.Context, user *User, clientIP, userAgent string) (*RefreshResponse, error) {
	return s.issueTokensForUser(ctx, user, clientIP, userAgent)
}

func (s *Service) LoginWithEmailAndPassword(ctx context.Context, email, password string) (*LoginResponse, error) {
	return s.LoginWithEmailAndPasswordWithContext(ctx, email, password, "", "")
}

// LoginWithEmailAndPasswordWithContext includes client context for security notifications
func (s *Service) LoginWithEmailAndPasswordWithContext(ctx context.Context, email, password, clientIP, userAgent string) (*LoginResponse, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	ip := strings.TrimSpace(clientIP)
	if ip == "" {
		ip = "unknown"
	}
	// Basic abuse protection for login attempts: best-effort, cache-backed.
	if !throttle.Allow(s.storage.cache, fmt.Sprintf("rl:auth:login:%s:%s", normalizedEmail, ip), 10*time.Minute, 20, 0) {
		return nil, ErrAuthRateLimited(nil)
	}

	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials(err)
	}
	if hash.ValidatePassword(password, user.Password) {
		tokens, err := s.issueTokensForUser(ctx, user, clientIP, userAgent)
		if err != nil {
			return nil, err
		}

		// Send login notification email asynchronously (best-effort)
		l := logger.FromContext(ctx)
		go func(u *User, ipAddr, ua string) {
			bg, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := s.Notification.SendLoginNotificationEmail(bg, u, ipAddr, ua); err != nil {
				l.Warn("failed to send login notification email", "error", err)
			}
		}(user, clientIP, userAgent)

		return ToLoginResponse(user, tokens.Token, tokens.RefreshToken), nil
	}
	return nil, ErrInvalidCredentials(nil)
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken, clientIP, userAgent string) (*RefreshResponse, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil, ErrInvalidOrExpiredToken(nil)
	}

	hash := auth.HashRefreshToken(refreshToken)
	sess, err := s.storage.GetSessionByTokenHash(ctx, hash, time.Now().UTC())
	if err != nil {
		if database.IsRecordNotFound(err) {
			return nil, ErrInvalidOrExpiredToken(err)
		}
		return nil, problem.InternalError().WithError(err)
	}

	user, err := s.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return nil, ErrInvalidOrExpiredToken(err)
	}

	// Rotate token: revoke old first to minimize replay window.
	if err := s.storage.RevokeSessionByTokenHash(ctx, hash); err != nil {
		return nil, problem.InternalError().WithError(err)
	}

	return s.issueTokensForUser(ctx, user, clientIP, userAgent)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return ErrInvalidOrExpiredToken(nil)
	}
	hash := auth.HashRefreshToken(refreshToken)
	if err := s.storage.RevokeSessionByTokenHash(ctx, hash); err != nil {
		return problem.InternalError().WithError(err)
	}
	return nil
}

func (s *Service) LogoutAll(ctx context.Context, refreshToken string) error {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return ErrInvalidOrExpiredToken(nil)
	}
	hash := auth.HashRefreshToken(refreshToken)
	sess, err := s.storage.GetSessionByTokenHash(ctx, hash, time.Now().UTC())
	if err != nil {
		if database.IsRecordNotFound(err) {
			return ErrInvalidOrExpiredToken(err)
		}
		return problem.InternalError().WithError(err)
	}

	user, err := s.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return ErrInvalidOrExpiredToken(err)
	}

	// Immediately invalidate all access tokens.
	user.AuthVersion++
	if err := s.storage.user.UpdateOne(ctx, user); err != nil {
		return problem.InternalError().WithError(err)
	}

	if err := s.storage.RevokeAllSessionsForUser(ctx, user.ID); err != nil {
		return problem.InternalError().WithError(err)
	}
	return nil
}

func (s *Service) LogoutOtherDevices(ctx context.Context, refreshToken string) error {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return ErrInvalidOrExpiredToken(nil)
	}
	hash := auth.HashRefreshToken(refreshToken)
	sess, err := s.storage.GetSessionByTokenHash(ctx, hash, time.Now().UTC())
	if err != nil {
		if database.IsRecordNotFound(err) {
			return ErrInvalidOrExpiredToken(err)
		}
		return problem.InternalError().WithError(err)
	}
	if err := s.storage.RevokeOtherSessionsForUser(ctx, sess.UserID, hash); err != nil {
		return problem.InternalError().WithError(err)
	}
	return nil
}

func (s *Service) CreateVerifyEmailToken(ctx context.Context, email string) (string, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	// Abuse protection: prevent spamming verification emails (best-effort, cache-backed).
	if !throttle.Allow(s.storage.cache, fmt.Sprintf("rl:auth:verify_email:%s", normalizedEmail), time.Hour, 5, 30*time.Second) {
		return "", ErrAuthRateLimited(nil)
	}

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
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	// Abuse protection: prevent spamming reset emails (best-effort, cache-backed).
	if !throttle.Allow(s.storage.cache, fmt.Sprintf("rl:auth:password_reset:%s", normalizedEmail), time.Hour, 5, 30*time.Second) {
		return "", ErrAuthRateLimited(nil)
	}

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
	user.AuthVersion++
	err = s.storage.user.UpdateOne(ctx, user)
	if err != nil {
		return problem.InternalError().WithError(err)
	}
	err = s.storage.ConsumePasswordResetToken(token)
	if err != nil {
		return problem.InternalError().WithError(err)
	}

	// Revoke all sessions so password resets kick everyone out.
	if err := s.storage.RevokeAllSessionsForUser(ctx, user.ID); err != nil {
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
	bootstrap := func(txCtx context.Context) (*User, *Workspace, error) {
		// Guard against existing user
		u, err := s.GetUserByEmail(txCtx, email)
		if err == nil && u != nil {
			return nil, nil, ErrInvalidCredentials(nil)
		}
		if err != nil && !database.IsRecordNotFound(err) {
			return nil, nil, err
		}

		ws := &Workspace{}
		if err := s.storage.workspace.CreateOne(txCtx, ws); err != nil {
			return nil, nil, err
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
			return nil, nil, err
		}
		ws.OwnerID = user.ID
		if stripeCustomerID != "" {
			ws.StripeCustomerID.String = stripeCustomerID
			ws.StripeCustomerID.Valid = true
		}
		if err := s.storage.workspace.UpdateOne(txCtx, ws); err != nil {
			return nil, nil, err
		}
		return user, ws, nil
	}

	// If we are already inside a DB transaction, reuse it to avoid nested transactions.
	// This is important for flows like onboarding completion that must be a single atomic commit.
	if ctx.Value(database.TxKey) != nil {
		return bootstrap(ctx)
	}

	var createdUser *User
	var createdWs *Workspace
	err := s.atomicProcessor.Exec(ctx, func(txCtx context.Context) error {
		u, ws, err := bootstrap(txCtx)
		if err != nil {
			return err
		}
		createdUser = u
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

	// Find the invitation (scoped to token payload)
	invitation, err := s.storage.invitation.FindOne(ctx,
		s.storage.invitation.ScopeID(payload.InvitationID),
		s.storage.invitation.ScopeEquals(UserInvitationSchema.WorkspaceID, payload.WorkspaceID),
		s.storage.invitation.ScopeEquals(UserInvitationSchema.Email, payload.Email),
	)
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

	// Find the invitation (scoped to token payload)
	invitation, err := s.storage.invitation.FindOne(ctx,
		s.storage.invitation.ScopeID(payload.InvitationID),
		s.storage.invitation.ScopeEquals(UserInvitationSchema.WorkspaceID, payload.WorkspaceID),
		s.storage.invitation.ScopeEquals(UserInvitationSchema.Email, payload.Email),
	)
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

	// Get target user scoped to the workspace (prevents ID probing).
	targetUser, err := s.GetWorkspaceUserByID(ctx, workspace.ID, targetUserID)
	if err != nil {
		return nil, err
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

// RevokeInvitation revokes a pending invitation.
// It enforces workspace scoping to prevent cross-tenant access (BOLA).
func (s *Service) RevokeInvitation(ctx context.Context, workspaceID string, invitationID string) error {
	invitation, err := s.storage.invitation.FindOne(
		ctx,
		s.storage.invitation.ScopeWorkspaceID(workspaceID),
		s.storage.invitation.ScopeID(invitationID),
	)
	if err != nil {
		return ErrInvitationNotFound(err)
	}

	if invitation.Status != InvitationStatusPending {
		return problem.Conflict("only pending invitations can be revoked")
	}

	invitation.Status = InvitationStatusRevoked
	return s.storage.invitation.UpdateOne(ctx, invitation)
}
