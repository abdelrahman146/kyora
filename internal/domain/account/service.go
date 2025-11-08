package account

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/email"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
)

type Service struct {
	bus             *bus.Bus
	storage         *Storage
	atomicProcessor atomic.AtomicProcessor
	notification    *Notification
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, bus *bus.Bus, emailClient email.Client) *Service {
	emailInfo := email.NewEmail()
	notification := NewNotification(emailClient, emailInfo)
	return &Service{
		storage:         storage,
		atomicProcessor: atomicProcessor,
		bus:             bus,
		notification:    notification,
	}
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	return s.storage.user.FindByID(ctx, id, s.storage.user.WithPreload(WorkspaceStruct))
}

func (s *Service) GetWorkspaceByID(ctx context.Context, id string) (*Workspace, error) {
	return s.storage.workspace.FindByID(ctx, id, s.storage.workspace.WithPreload(UserStruct))
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
			if err := s.notification.SendLoginNotificationEmail(context.Background(), user, clientIP, userAgent); err != nil {
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

	err = s.notification.SendEmailVerificationEmail(ctx, user, token, expAt)
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
	err = s.notification.SendForgotPasswordEmail(ctx, user, token, expAt)
	if err != nil {
		// Log error but don't fail the token creation
		// Fall back to the bus system
		s.bus.Emit(bus.ResetPasswordTopic, &bus.ResetPasswordEvent{
			Email:    user.Email,
			ExpireAt: expAt,
			Token:    token,
		})
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
	err = s.notification.SendPasswordResetConfirmationEmail(ctx, user, "")
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
