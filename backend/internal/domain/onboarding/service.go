package onboarding

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/email"
	"github.com/abdelrahman146/kyora/internal/platform/logger"
	pactomic "github.com/abdelrahman146/kyora/internal/platform/types/atomic"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
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
	}
	// Try to fetch plan and determine paid/free
	plan, err := s.billing.GetPlanByDescriptor(ctx, planDescriptor)
	if err != nil || plan == nil {
		return nil, ErrPlanNotFound(err)
	}
	isPaid := !plan.Price.IsZero()
	// Resume if active session exists
	if existing, err := s.storage.GetActiveByEmail(ctx, emailStr); err == nil && existing != nil {
		// Update plan if changed
		existing.PlanID = plan.ID
		existing.PlanDescriptor = plan.Descriptor
		existing.IsPaidPlan = isPaid
		if existing.Stage == "" {
			existing.Stage = StagePlanSelected
		}
		_ = s.storage.UpdateSession(ctx, existing)
		return existing, nil
	}
	// Clear expired sessions to avoid unique constraint conflicts
	_ = s.storage.DeleteExpiredSessionsByEmail(ctx, emailStr)
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
func (s *Service) GenerateEmailOTP(ctx context.Context, token string) error {
	sess, err := s.storage.GetByToken(ctx, token)
	if err != nil || sess == nil {
		return ErrSessionNotFound(err)
	}
	if time.Now().After(sess.ExpiresAt) {
		return ErrSessionExpired(nil)
	}
	if sess.Stage != StagePlanSelected && sess.Stage != StageIdentityPending {
		return ErrInvalidStage(nil, string(StagePlanSelected))
	}
	// rate limit: max 5 per hour and at least 30s between requests
	if err := s.throttle("rl:otp:"+sess.Email, time.Hour, 5, 30*time.Second); err != nil {
		return err
	}
	code, err := id.RandomNumber(6)
	if err != nil {
		return err
	}
	hashed, err := hash.Password(code)
	if err != nil {
		return err
	}
	exp := time.Now().Add(15 * time.Minute)
	sess.OTPHash = hashed
	sess.OTPExpiry = &exp
	sess.Stage = StageIdentityPending
	if err := s.storage.UpdateSession(ctx, sess); err != nil {
		return err
	}
	// Send OTP code email for onboarding. This must include the OTP code.
	// We intentionally avoid linking to the account /verify-email flow because onboarding uses a different API.
	if s.emailClient == nil {
		return fmt.Errorf("email client not available")
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
		return err
	}
	logger.FromContext(ctx).Info("onboarding email OTP generated", "session", sess.ID)
	return nil
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
	if sess.OTPExpiry == nil || time.Now().After(*sess.OTPExpiry) {
		return nil, ErrInvalidOTP(nil)
	}
	if !hash.ValidatePassword(code, sess.OTPHash) {
		return nil, ErrInvalidOTP(nil)
	}
	passHash, err := hash.Password(password)
	if err != nil {
		return nil, err
	}
	sess.EmailVerified = true
	sess.FirstName = firstName
	sess.LastName = lastName
	sess.PasswordHash = passHash
	sess.Stage = StageIdentityVerified
	if err := s.storage.UpdateSession(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
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
	if err := s.throttle("rl:pay:"+sess.ID, 10*time.Minute, 3, 30*time.Second); err != nil {
		return "", err
	}
	plan, err := s.billing.GetPlanByDescriptor(ctx, sess.PlanDescriptor)
	if err != nil || plan == nil {
		return "", ErrPlanNotFound(err)
	}
	if plan.StripePlanID == "" {
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
		LineItems:  []*stripelib.CheckoutSessionLineItemParams{{Price: stripelib.String(plan.StripePlanID), Quantity: stripelib.Int64(1)}},
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
func (s *Service) CompleteOnboarding(ctx context.Context, token string) (user *account.User, jwt string, err error) {
	sess, err := s.storage.GetByToken(ctx, token)
	if err != nil || sess == nil {
		return nil, "", ErrSessionNotFound(err)
	}
	if time.Now().After(sess.ExpiresAt) {
		return nil, "", ErrSessionExpired(nil)
	}
	if sess.Stage != StageReadyToCommit {
		return nil, "", ErrInvalidStage(nil, string(StageReadyToCommit))
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
		return nil
	}, pactomic.WithRetries(2))
	if err != nil {
		return nil, "", err
	}
	// Send welcome email (best effort)
	go func() { _ = s.emailNotif.SendWelcomeEmail(context.Background(), createdUser) }()
	// Invalidate session by deleting it
	now := time.Now()
	sess.CommittedAt = &now
	sess.Stage = StageCommitted
	_ = s.storage.UpdateSession(context.Background(), sess)
	// Generate JWT
	jwtToken, jwtErr := auth.NewJwtToken(createdUser.ID, createdUser.WorkspaceID)
	if jwtErr != nil {
		return createdUser, "", jwtErr
	}
	return createdUser, jwtToken, nil
}

// throttle implements a simple token bucket using cache with JSON state
type throttleState struct {
	Count int   `json:"count"`
	Last  int64 `json:"last"` // unix seconds
}

func (s *Service) throttle(key string, window time.Duration, max int, minInterval time.Duration) error {
	if s.storage.cache == nil {
		return nil
	}
	now := time.Now()
	// read state
	var st throttleState
	if data, err := s.storage.cache.Get(key); err == nil && len(data) > 0 {
		_ = s.storage.cache.Unmarshal(data, &st)
	}
	// enforce min interval
	if st.Last != 0 && now.Sub(time.Unix(st.Last, 0)) < minInterval {
		return ErrRateLimited(nil)
	}
	// increment and cap
	st.Count++
	st.Last = now.Unix()
	if st.Count > max {
		// write back to ensure TTL stays
		if b, err := s.storage.cache.Marshal(st); err == nil {
			_ = s.storage.cache.SetX(key, b, int32(window.Seconds()))
		}
		return ErrRateLimited(nil)
	}
	// persist with TTL window
	if b, err := s.storage.cache.Marshal(st); err == nil {
		_ = s.storage.cache.SetX(key, b, int32(window.Seconds()))
	}
	return nil
}

// Expose selected helpers for other packages (webhooks)
// removed EnsurePlanSynced wrapper (not needed externally yet)
