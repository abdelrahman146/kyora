package onboarding

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/email"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	pactomic "github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/abdelrahman146/kyora/internal/platform/utils/throttle"
	stripelib "github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/checkout/session"
	"github.com/stripe/stripe-go/v83/customer"
)

type Service struct {
	storage         *Storage
	atomicProcessor pactomic.AtomicProcessor
	account         *account.Service
	billing         *billing.Service
	emailNotif      *account.Notification
	emailClient     email.Client
	emailInfo       email.EmailInfo
	business        *business.Service
}

func NewService(storage *Storage, atomic pactomic.AtomicProcessor, accountSvc *account.Service, billingSvc *billing.Service, businessSvc *business.Service, emailClient email.Client) *Service {
	emailInfo := email.NewEmail()
	notif := account.NewNotification(emailClient, emailInfo)
	return &Service{
		storage:         storage,
		atomicProcessor: atomic,
		account:         accountSvc,
		billing:         billingSvc,
		business:        businessSvc,
		emailNotif:      notif,
		emailClient:     emailClient,
		emailInfo:       emailInfo,
	}
}

// StartSession initializes or resumes a session for an email and plan.
func (s *Service) StartSession(ctx context.Context, emailStr, planDescriptor string, expiresIn time.Duration) (*OnboardingSession, error) {
	// Ensure email is not already a completed user
	if _, err := s.account.GetUserByEmail(ctx, emailStr); err == nil {
		return nil, ErrEmailAlreadyExists(nil)
	} else if err != nil && !database.IsRecordNotFound(err) {
		return nil, err
	}
	// Try to fetch plan and determine paid/free
	plan, err := s.billing.GetPlanByDescriptor(ctx, planDescriptor)
	if err != nil || plan == nil {
		return nil, ErrPlanNotFound(err)
	}
	isPaid := !plan.Price.IsZero()
	// Resume if active session exists
	existing, err := s.storage.GetActiveByEmail(ctx, emailStr)
	if err == nil {
		// Update plan if changed
		existing.PlanID = plan.ID
		existing.PlanDescriptor = plan.Descriptor
		existing.IsPaidPlan = isPaid
		if existing.Stage == "" {
			existing.Stage = StagePlanSelected
		}
		if err := s.storage.UpdateSession(ctx, existing); err != nil {
			return nil, ErrSessionUpdateFailed(err)
		}
		return existing, nil
	}
	if err != nil && !database.IsRecordNotFound(err) {
		return nil, err
	}
	// Clear expired sessions to avoid unique constraint conflicts
	if err := s.storage.DeleteExpiredSessionsByEmail(ctx, emailStr); err != nil {
		return nil, ErrSessionCleanupFailed(err)
	}
	sess := &OnboardingSession{
		Email:          emailStr,
		PlanID:         plan.ID,
		PlanDescriptor: plan.Descriptor,
		IsPaidPlan:     isPaid,
		Method:         IdentityEmail,
		Stage:          StagePlanSelected,
		PaymentStatus: func() PaymentStatus {
			if isPaid {
				return PaymentStatusPending
			}
			return PaymentStatusSkipped
		}(),
		ExpiresAt: time.Now().Add(expiresIn),
	}
	if err := s.storage.CreateSession(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

// GenerateEmailOTP issues a 6-digit code and emails it. Stores hash + expiry.
func (s *Service) GenerateEmailOTP(ctx context.Context, token string) (time.Duration, error) {
	sess, err := s.storage.GetByToken(ctx, token)
	if err != nil || sess == nil {
		return 0, ErrSessionNotFound(err)
	}
	if time.Now().After(sess.ExpiresAt) {
		return 0, ErrSessionExpired(nil)
	}
	if sess.Stage != StagePlanSelected && sess.Stage != StageIdentityPending {
		return 0, ErrInvalidStage(nil, string(StagePlanSelected))
	}
	// Cooldown strategy: per email, allow one OTP request then block resends for 2 minutes.
	// This is intentionally simpler and more predictable than token-bucket throttling.
	cooldown := 2 * time.Minute
	cooldownKey := "cd:onboarding:otp:" + strings.ToLower(strings.TrimSpace(sess.Email))
	allowed, retryAfter := throttle.Cooldown(s.storage.cache, cooldownKey, cooldown)
	if !allowed {
		return 0, ErrRateLimitedRetryAfter(nil, retryAfter)
	}
	code, err := id.RandomNumber(6)
	if err != nil {
		return 0, err
	}
	hashed, err := hash.Password(code)
	if err != nil {
		return 0, err
	}
	exp := time.Now().Add(15 * time.Minute)
	sess.OTPHash = hashed
	sess.OTPExpiry = &exp
	sess.Stage = StageIdentityPending
	if err := s.storage.UpdateSession(ctx, sess); err != nil {
		return 0, err
	}
	// Send OTP code email for onboarding. This must include the OTP code.
	// We intentionally avoid linking to the account /verify-email flow because onboarding uses a different API.
	if s.emailClient == nil {
		if s.storage.cache != nil {
			_ = s.storage.cache.Delete(cooldownKey)
		}
		return 0, fmt.Errorf("email client not available")
	}
	userName := strings.TrimSpace(strings.Join([]string{sess.FirstName, sess.LastName}, " "))
	if userName == "" {
		parts := strings.Split(sess.Email, "@")
		if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
			userName = parts[0]
		} else {
			userName = "there"
		}
	}
	from := s.emailInfo.FormattedFrom()
	_, err = s.emailClient.SendTemplate(ctx, email.TemplateOnboardingEmailOTP, []string{sess.Email}, from, "", map[string]any{
		"userName":     userName,
		"otpCode":      code,
		"expiryTime":   "15 minutes",
		"productName":  s.emailInfo.ProductName,
		"supportEmail": s.emailInfo.SupportEmail,
		"helpURL":      s.emailInfo.HelpURL,
		"currentYear":  fmt.Sprintf("%d", time.Now().Year()),
	})
	if err != nil {
		if s.storage.cache != nil {
			_ = s.storage.cache.Delete(cooldownKey)
		}
		return 0, err
	}
	logger.FromContext(ctx).Info("onboarding email OTP generated", "session", sess.ID)
	return cooldown, nil
}

// VerifyEmailOTP validates code and stores password+profile
func (s *Service) VerifyEmailOTP(ctx context.Context, token, code, firstName, lastName, password string) (*OnboardingSession, error) {
	sess, err := s.storage.GetByToken(ctx, token)
	if err != nil || sess == nil {
		return nil, ErrSessionNotFound(err)
	}
	if time.Now().After(sess.ExpiresAt) {
		return nil, ErrSessionExpired(nil)
	}
	if sess.Stage != StageIdentityPending {
		return nil, ErrInvalidStage(nil, string(StageIdentityPending))
	}
	if sess.OTPExpiry == nil || time.Now().After(*sess.OTPExpiry) {
		return nil, ErrInvalidOTP(nil)
	}
	if !hash.ValidatePassword(code, sess.OTPHash) {
		return nil, ErrInvalidOTP(nil)
	}
	nFirst, err := normalizeHumanName("firstName", firstName)
	if err != nil {
		return nil, err
	}
	nLast, err := normalizeHumanName("lastName", lastName)
	if err != nil {
		return nil, err
	}
	if password == "" {
		return nil, schemaValidationError("password", "is required")
	}
	if len([]rune(password)) < 8 {
		return nil, schemaValidationError("password", "must be at least 8 characters")
	}
	passHash, err := hash.Password(password)
	if err != nil {
		return nil, err
	}
	sess.EmailVerified = true
	sess.FirstName = nFirst
	sess.LastName = nLast
	sess.PasswordHash = passHash
	sess.Stage = StageIdentityVerified
	if err := s.storage.UpdateSession(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

var humanNameRegex = regexp.MustCompile(`^[\p{L}\p{M}][\p{L}\p{M} '\-]{0,99}$`)

func normalizeHumanName(field, v string) (string, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", schemaValidationError(field, "is required")
	}
	// Disallow NUL/control characters. Postgres TEXT cannot store NUL bytes.
	for _, r := range v {
		if r == 0 || r < 32 || r == 127 {
			return "", schemaValidationError(field, "contains invalid characters")
		}
	}
	if len([]rune(v)) > 100 {
		return "", schemaValidationError(field, "is too long")
	}
	if !humanNameRegex.MatchString(v) {
		return "", schemaValidationError(field, "contains invalid characters")
	}
	return v, nil
}

func schemaValidationError(field, msg string) error {
	return problem.BadRequest(msg).With("field", field).WithCode("onboarding.validation_error")
}

// SetOAuthIdentity stages identity from Google and generates a random password
func (s *Service) SetOAuthIdentity(ctx context.Context, token string, googleCode string) (*OnboardingSession, error) {
	sess, err := s.storage.GetByToken(ctx, token)
	if err != nil || sess == nil {
		return nil, ErrSessionNotFound(err)
	}
	if time.Now().After(sess.ExpiresAt) {
		return nil, ErrSessionExpired(nil)
	}
	info, err := s.account.ExchangeGoogleCodeAndFetchUser(ctx, googleCode)
	if err != nil {
		return nil, err
	}
	randPass, _ := id.RandomString(24)
	passHash, err := hash.Password(randPass)
	if err != nil {
		return nil, err
	}
	sess.Method = IdentityGoogle
	sess.Email = info.Email
	sess.EmailVerified = true
	sess.FirstName = info.GivenName
	sess.LastName = info.FamilyName
	sess.PasswordHash = passHash
	sess.Stage = StageIdentityVerified
	if err := s.storage.UpdateSession(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

// SetBusinessDetails stages business data
func (s *Service) SetBusinessDetails(ctx context.Context, token string, name, descriptor, country, currency string) (*OnboardingSession, error) {
	sess, err := s.storage.GetByToken(ctx, token)
	if err != nil || sess == nil {
		return nil, ErrSessionNotFound(err)
	}
	if time.Now().After(sess.ExpiresAt) {
		return nil, ErrSessionExpired(nil)
	}
	if sess.Stage != StageIdentityVerified && sess.Stage != StageBusinessStaged {
		return nil, ErrInvalidStage(nil, string(StageIdentityVerified))
	}
	sess.BusinessName = name
	sess.BusinessDescriptor = descriptor
	sess.BusinessCountry = country
	sess.BusinessCurrency = currency
	sess.Stage = StageBusinessStaged
	if !sess.IsPaidPlan {
		sess.Stage = StageReadyToCommit
		sess.PaymentStatus = PaymentStatusSkipped
	} else {
		sess.Stage = StagePaymentPending
		sess.PaymentStatus = PaymentStatusPending
	}
	if err := s.storage.UpdateSession(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

// InitiatePayment creates Stripe customer and checkout session when plan is paid. Returns URL.
func (s *Service) InitiatePayment(ctx context.Context, token string, successURL, cancelURL string) (string, error) {
	sess, err := s.storage.GetByToken(ctx, token)
	if err != nil || sess == nil {
		return "", ErrSessionNotFound(err)
	}
	if !sess.IsPaidPlan {
		return "", nil
	}
	// Throttle duplicate attempts: max 3 per 10 minutes, at least 30s apart
	if !throttle.Allow(s.storage.cache, "rl:pay:"+sess.ID, 10*time.Minute, 3, 30*time.Second) {
		return "", ErrRateLimited(nil)
	}
	plan, err := s.billing.GetPlanByDescriptor(ctx, sess.PlanDescriptor)
	if err != nil || plan == nil {
		return "", ErrPlanNotFound(err)
	}
	if plan.StripePlanID == nil || *plan.StripePlanID == "" {
		return "", ErrStripeOperation(fmt.Errorf("plan not available for checkout"))
	}
	// If we already have a pending checkout session, reuse its URL
	if sess.CheckoutSessionID != "" && sess.PaymentStatus == PaymentStatusPending {
		if cs, err := session.Get(sess.CheckoutSessionID, nil); err == nil && cs != nil && cs.URL != "" {
			return cs.URL, nil
		}
	}
	// Create or reuse Stripe customer with email
	if sess.StripeCustomerID == "" {
		c, err := customer.New(&stripelib.CustomerParams{Email: stripelib.String(sess.Email), Metadata: map[string]string{"onboarding_session_id": sess.ID}})
		if err != nil {
			return "", ErrStripeOperation(err)
		}
		sess.StripeCustomerID = c.ID
	}
	params := &stripelib.CheckoutSessionParams{
		Customer:   stripelib.String(sess.StripeCustomerID),
		Mode:       stripelib.String(string(stripelib.CheckoutSessionModeSubscription)),
		LineItems:  []*stripelib.CheckoutSessionLineItemParams{{Price: stripelib.String(*plan.StripePlanID), Quantity: stripelib.Int64(1)}},
		SuccessURL: stripelib.String(successURL),
		CancelURL:  stripelib.String(cancelURL),
		Metadata:   map[string]string{"onboarding_session_id": sess.ID, "plan_id": plan.ID},
	}
	sessParamsID := fmt.Sprintf("onboarding_checkout_%s", sess.ID)
	params.SetIdempotencyKey(sessParamsID)
	cs, err := session.New(params)
	if err != nil {
		return "", ErrStripeOperation(err)
	}
	sess.CheckoutSessionID = cs.ID
	if err := s.storage.UpdateSession(ctx, sess); err != nil {
		return "", err
	}
	return cs.URL, nil
}

// MarkPaymentSucceeded flags when webhook confirms checkout.session.completed
func (s *Service) MarkPaymentSucceeded(ctx context.Context, sessionID, stripeSubID string) error {
	// lookup by checkout session id
	rec, err := s.storage.session.FindOne(ctx, s.storage.session.ScopeEquals(schema.NewField("checkout_session_id", "checkoutSessionId"), sessionID))
	if err != nil || rec == nil {
		return ErrSessionNotFound(err)
	}
	rec.PaymentStatus = PaymentStatusSucceeded
	rec.StripeSubID = stripeSubID
	rec.Stage = StageReadyToCommit
	return s.storage.UpdateSession(ctx, rec)
}

// CompleteOnboarding performs a single transactional commit into permanent tables.
func (s *Service) CompleteOnboarding(ctx context.Context, token, clientIP, userAgent string) (user *account.User, jwt string, refreshToken string, err error) {
	sess, err := s.storage.GetByToken(ctx, token)
	if err != nil || sess == nil {
		return nil, "", "", ErrSessionNotFound(err)
	}
	if time.Now().After(sess.ExpiresAt) {
		return nil, "", "", ErrSessionExpired(nil)
	}
	if sess.Stage != StageReadyToCommit {
		return nil, "", "", ErrInvalidStage(nil, string(StageReadyToCommit))
	}
	var createdUser *account.User
	var createdWorkspace *account.Workspace
	err = s.atomicProcessor.Exec(ctx, func(txCtx context.Context) error {
		// bootstrap workspace & owner via account service abstraction
		u, ws, err := s.account.BootstrapWorkspaceAndOwner(txCtx, sess.FirstName, sess.LastName, sess.Email, sess.PasswordHash, sess.EmailVerified, sess.StripeCustomerID)
		if err != nil {
			return err
		}
		createdUser = u
		createdWorkspace = ws
		// create business using business service (actor = owner)
		if s.business != nil {
			_, err = s.business.CreateBusiness(txCtx, createdUser, &business.CreateBusinessInput{Descriptor: sess.BusinessDescriptor, Name: sess.BusinessName, CountryCode: sess.BusinessCountry, Currency: sess.BusinessCurrency})
			if err != nil {
				return err
			}
		} else {
			// fallback direct repository if service missing
			bizRec := &business.Business{WorkspaceID: ws.ID, Name: sess.BusinessName, Descriptor: sess.BusinessDescriptor, CountryCode: sess.BusinessCountry, Currency: sess.BusinessCurrency}
			if err := s.storage.CreateBusiness(txCtx, bizRec); err != nil {
				return err
			}
		}
		if sess.IsPaidPlan {
			plan, err := s.billing.GetPlanByID(txCtx, sess.PlanID)
			if err != nil {
				return err
			}
			if _, err := s.billing.CreateOrUpdateSubscription(txCtx, createdWorkspace, plan); err != nil {
				return err
			}
		}

		now := time.Now()
		sess.CommittedAt = &now
		sess.Stage = StageCommitted
		if err := s.storage.UpdateSession(txCtx, sess); err != nil {
			return err
		}

		return nil
	}, pactomic.WithRetries(2))
	if err != nil {
		return nil, "", "", err
	}
	// Send welcome email (best effort)
	l := logger.FromContext(ctx)
	go func(u *account.User) {
		bg, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.emailNotif.SendWelcomeEmail(bg, u); err != nil {
			l.Warn("failed to send welcome email", "error", err)
		}
	}(createdUser)
	// Generate JWT
	tokens, jwtErr := s.account.IssueTokensForUserWithContext(ctx, createdUser, clientIP, userAgent)
	if jwtErr != nil {
		return createdUser, "", "", jwtErr
	}
	return createdUser, tokens.Token, tokens.RefreshToken, nil
}

// Expose selected helpers for other packages (webhooks)
// removed EnsurePlanSynced wrapper (not needed externally yet)

// GetSession retrieves an active onboarding session by token.
// Returns error if session not found, expired, or already committed.
func (s *Service) GetSession(ctx context.Context, sessionToken string) (*OnboardingSession, error) {
	sess, err := s.storage.GetByToken(ctx, sessionToken)
	if err != nil {
		if database.IsRecordNotFound(err) {
			return nil, ErrSessionNotFound(err)
		}
		return nil, err
	}

	// Check expiry
	if time.Now().After(sess.ExpiresAt) {
		return nil, ErrSessionExpired(nil)
	}

	// Check if already committed
	if sess.Stage == StageCommitted || sess.CommittedAt != nil {
		return nil, ErrSessionExpired(nil)
	}

	return sess, nil
}

// DeleteSession deletes an onboarding session by token.
// Allows users to cancel/restart onboarding flow.
func (s *Service) DeleteSession(ctx context.Context, sessionToken string) error {
	sess, err := s.storage.GetByToken(ctx, sessionToken)
	if err != nil {
		if database.IsRecordNotFound(err) {
			return ErrSessionNotFound(err)
		}
		return err
	}

	// Don't allow deleting already committed sessions
	if sess.Stage == StageCommitted || sess.CommittedAt != nil {
		return ErrSessionAlreadyCommitted(nil)
	}

	return s.storage.DeleteSession(ctx, sess)
}
