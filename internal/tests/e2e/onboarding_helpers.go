package e2e

import (
	"context"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/domain/onboarding"
	"github.com/abdelrahman146/kyora/internal/platform/cache"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
)

// OnboardingTestHelper provides reusable helpers for onboarding tests
type OnboardingTestHelper struct {
	db                *database.Database
	client            *testutils.HTTPClient
	onboardingStorage *onboarding.Storage
	accountStorage    *account.Storage
	billingStorage    *billing.Storage
}

// NewOnboardingTestHelper creates a new onboarding test helper
func NewOnboardingTestHelper(db *database.Database, baseURL string) *OnboardingTestHelper {
	// Initialize cache for storage layers
	cacheClient := cache.NewConnection([]string{"localhost:11211"})

	return &OnboardingTestHelper{
		db:                db,
		client:            testutils.NewHTTPClient(baseURL),
		onboardingStorage: onboarding.NewStorage(db, cacheClient),
		accountStorage:    account.NewStorage(db, cacheClient),
		billingStorage:    billing.NewStorage(db, cacheClient),
	}
}

// EnsureTestPlan creates a plan if it doesn't exist with proper features and limits
func (h *OnboardingTestHelper) EnsureTestPlan(descriptor, name string, price float64) {
	// Check if plan already exists
	var existingPlan billing.Plan
	err := h.db.GetDB().Where("descriptor = ?", descriptor).First(&existingPlan).Error
	if err == nil {
		// Plan already exists
		return
	}

	// Create features
	features := billing.PlanFeature{
		CustomerManagement:       true,
		InventoryManagement:      true,
		OrderManagement:          true,
		ExpenseManagement:        true,
		Accounting:               true,
		BasicAnalytics:           true,
		FinancialReports:         true,
		DataImport:               price > 0,
		DataExport:               price > 0,
		AdvancedAnalytics:        price > 0,
		AdvancedFinancialReports: price > 0,
		OrderPaymentLinks:        price > 0,
		InvoiceGeneration:        price > 0,
		ExportAnalyticsData:      false,
		AIBusinessAssistant:      false,
	}

	// Create limits
	maxOrders := int64(25)
	maxTeamMembers := int64(1)
	maxBusinesses := int64(1)
	if price > 0 {
		maxOrders = 500
		maxTeamMembers = 5
		maxBusinesses = 3
	}
	limits := billing.PlanLimit{
		MaxOrdersPerMonth: maxOrders,
		MaxTeamMembers:    maxTeamMembers,
		MaxBusinesses:     maxBusinesses,
	}

	description := "Test plan for " + name

	plan := billing.Plan{
		ID:           "plan_" + descriptor,
		Descriptor:   descriptor,
		Name:         name,
		Description:  description,
		Price:        decimal.NewFromFloat(price),
		Currency:     "aed",
		StripePlanID: "",
		BillingCycle: billing.BillingCycleMonthly,
		Features:     features,
		Limits:       limits,
	}

	// Use direct database access for test plan creation
	// Production plans are managed by billing storage init, but test plans need manual creation
	h.db.GetDB().Create(&plan)
}

// CreateTestUser creates a user and workspace for testing email already registered scenarios
func (h *OnboardingTestHelper) CreateTestUser(email, password, firstName, lastName string) error {
	ctx := context.Background()
	workspaceID := "ws_test_" + email
	userID := "user_test_" + email

	// Hash password
	hashedPass, err := hash.Password(password)
	if err != nil {
		return err
	}

	// Create workspace using onboarding storage (which has access to workspace repo)
	workspace := &account.Workspace{
		ID:      workspaceID,
		OwnerID: userID,
	}
	if err := h.onboardingStorage.CreateWorkspace(ctx, workspace); err != nil {
		return err
	}

	// Create user using onboarding storage (which has access to user repo)
	user := &account.User{
		ID:              userID,
		WorkspaceID:     workspaceID,
		Email:           email,
		Password:        hashedPass,
		FirstName:       firstName,
		LastName:        lastName,
		IsEmailVerified: true,
		Role:            "admin",
	}
	return h.onboardingStorage.CreateUser(ctx, user)
} // CreateOnboardingSession creates a session via API and returns token
func (h *OnboardingTestHelper) CreateOnboardingSession(email, planDescriptor string) (string, error) {
	h.EnsureTestPlan(planDescriptor, "Test Plan", 0.0)

	payload := map[string]interface{}{
		"email":          email,
		"planDescriptor": planDescriptor,
	}
	resp, err := h.client.Post("/api/onboarding/start", payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := testutils.DecodeJSON(resp, &result); err != nil {
		return "", err
	}

	// Check if session token exists (API might have returned error)
	if result["sessionToken"] == nil {
		return "", nil
	}

	return result["sessionToken"].(string), nil
}

// SetSessionOTP sets OTP for a session
func (h *OnboardingTestHelper) SetSessionOTP(token, otp string, expiresIn time.Duration) error {
	ctx := context.Background()

	// Get session using storage
	sess, err := h.onboardingStorage.GetByToken(ctx, token)
	if err != nil {
		return err
	}

	// Hash OTP
	hashedOTP, err := hash.Password(otp)
	if err != nil {
		return err
	}

	// Update session fields
	expiry := time.Now().Add(expiresIn)
	sess.OTPHash = hashedOTP
	sess.OTPExpiry = &expiry
	sess.Stage = onboarding.StageIdentityPending

	// Update using storage
	return h.onboardingStorage.UpdateSession(ctx, sess)
}

// SetExpiredOTP sets an expired OTP for testing
func (h *OnboardingTestHelper) SetExpiredOTP(token string) error {
	ctx := context.Background()

	// Get session using storage
	sess, err := h.onboardingStorage.GetByToken(ctx, token)
	if err != nil {
		return err
	}

	// Set expired time
	expiredTime := time.Now().UTC().Add(-20 * time.Minute)
	sess.OTPExpiry = &expiredTime

	// Update using storage
	return h.onboardingStorage.UpdateSession(ctx, sess)
} // UpdateSessionStage updates session stage
func (h *OnboardingTestHelper) UpdateSessionStage(token, stage string) error {
	ctx := context.Background()

	// Get session using storage
	sess, err := h.onboardingStorage.GetByToken(ctx, token)
	if err != nil {
		return err
	}

	// Update stage (convert string to SessionStage type)
	sess.Stage = onboarding.SessionStage(stage)

	// Update using storage
	return h.onboardingStorage.UpdateSession(ctx, sess)
}

// ExpireSession sets session to expired
func (h *OnboardingTestHelper) ExpireSession(token string) error {
	ctx := context.Background()

	// Get session using storage
	sess, err := h.onboardingStorage.GetByToken(ctx, token)
	if err != nil {
		return err
	}

	// Set expired time in UTC (1 hour in the past)
	expiredTime := time.Now().UTC().Add(-1 * time.Hour)
	sess.ExpiresAt = expiredTime

	// Update using storage
	return h.onboardingStorage.UpdateSession(ctx, sess)
}

// GetSession retrieves session data
func (h *OnboardingTestHelper) GetSession(token string) (map[string]interface{}, error) {
	ctx := context.Background()

	// Get session using storage
	sess, err := h.onboardingStorage.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Convert to map for compatibility with existing tests
	result := make(map[string]interface{})
	result["stage"] = string(sess.Stage)
	result["email"] = sess.Email
	result["emailVerified"] = sess.EmailVerified
	result["otpExpiry"] = sess.OTPExpiry
	return result, nil
}

// CreateSessionWithOTP creates session with OTP set
func (h *OnboardingTestHelper) CreateSessionWithOTP(email, planDescriptor, otp string) (string, error) {
	token, err := h.CreateOnboardingSession(email, planDescriptor)
	if err != nil || token == "" {
		return "", err
	}

	if err := h.SetSessionOTP(token, otp, 15*time.Minute); err != nil {
		return "", err
	}

	return token, nil
}

// CreateVerifiedSession creates a session at identity_verified stage
func (h *OnboardingTestHelper) CreateVerifiedSession(email, planDescriptor string) (string, error) {
	ctx := context.Background()

	token, err := h.CreateOnboardingSession(email, planDescriptor)
	if err != nil || token == "" {
		return "", err
	}

	// Get session using storage
	sess, err := h.onboardingStorage.GetByToken(ctx, token)
	if err != nil {
		return "", err
	}

	// Hash password
	hashedPass, err := hash.Password("TestPassword123!")
	if err != nil {
		return "", err
	}

	// Update session fields
	sess.Stage = onboarding.StageIdentityVerified
	sess.EmailVerified = true
	sess.FirstName = "Test"
	sess.LastName = "User"
	sess.PasswordHash = hashedPass
	sess.Method = onboarding.IdentityEmail

	// Update using storage
	err = h.onboardingStorage.UpdateSession(ctx, sess)
	return token, err
} // CreateBusinessStagedSession creates session at ready_to_commit stage
func (h *OnboardingTestHelper) CreateBusinessStagedSession(email, planDescriptor string) (string, error) {
	ctx := context.Background()

	token, err := h.CreateVerifiedSession(email, planDescriptor)
	if err != nil || token == "" {
		return "", err
	}

	// Get session using storage
	sess, err := h.onboardingStorage.GetByToken(ctx, token)
	if err != nil {
		return "", err
	}

	// Update session fields for business stage
	sess.Stage = onboarding.StageReadyToCommit
	sess.BusinessName = "Test Business"
	sess.BusinessDescriptor = "test-business"
	sess.BusinessCountry = "AE"
	sess.BusinessCurrency = "AED"
	sess.PaymentStatus = onboarding.PaymentStatusSkipped

	// Update using storage
	err = h.onboardingStorage.UpdateSession(ctx, sess)
	return token, err
}
