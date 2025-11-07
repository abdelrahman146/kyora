package account

import (
	"context"

	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/bus"
	"github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
)

type Service struct {
	bus              *bus.Bus
	storage          *Storage
	atomicProcessor  atomic.AtomicProcessor
	emailIntegration *EmailIntegration
}

func NewService(storage *Storage, atomicProcessor atomic.AtomicProcessor, bus *bus.Bus, emailIntegration *EmailIntegration) *Service {
	return &Service{
		storage:          storage,
		atomicProcessor:  atomicProcessor,
		bus:              bus,
		emailIntegration: emailIntegration,
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
		if s.emailIntegration != nil {
			go func() {
				if err := s.emailIntegration.SendLoginNotificationEmail(context.Background(), user, clientIP, userAgent); err != nil {
					// Log error but don't fail the login process
					// Login notifications are a security feature but shouldn't block user access
				}
			}()
		}

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

	// Try to send email using the email integration
	if s.emailIntegration != nil {
		err = s.emailIntegration.SendEmailVerificationEmail(ctx, user, token, expAt)
		if err != nil {
			// Log error but don't fail the token creation
			// Fall back to the bus system
			s.bus.Emit(bus.VerifyEmailTopic, &bus.VerifyEmailEvent{
				Email:    user.Email,
				ExpireAt: expAt,
				Token:    token,
			})
		}
	} else {
		// Fallback to old bus system if email integration not available
		s.bus.Emit(bus.VerifyEmailTopic, &bus.VerifyEmailEvent{
			Email:    user.Email,
			ExpireAt: expAt,
			Token:    token,
		})
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
	if s.emailIntegration != nil {
		err = s.emailIntegration.SendForgotPasswordEmail(ctx, user, token, expAt)
		if err != nil {
			// Log error but don't fail the token creation
			// Fall back to the bus system
			s.bus.Emit(bus.ResetPasswordTopic, &bus.ResetPasswordEvent{
				Email:    user.Email,
				ExpireAt: expAt,
				Token:    token,
			})
		}
	} else {
		// Fallback to old bus system if email integration not available
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
	if s.emailIntegration != nil {
		// We don't have client IP in this context, so pass empty string
		// In a real implementation, you'd extract this from the HTTP request
		err = s.emailIntegration.SendPasswordResetConfirmationEmail(ctx, user, "")
		if err != nil {
			// Log error but don't fail the password reset operation
			// The password has already been changed successfully
		}
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
