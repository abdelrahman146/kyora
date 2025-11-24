package e2e

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/shopspring/decimal"
)

// OnboardingTestHelper provides reusable helpers for onboarding tests
type OnboardingTestHelper struct {
	db     *database.Database
	client *testutils.HTTPClient
}

// NewOnboardingTestHelper creates a new onboarding test helper
func NewOnboardingTestHelper(db *database.Database, baseURL string) *OnboardingTestHelper {
	return &OnboardingTestHelper{
		db:     db,
		client: testutils.NewHTTPClient(baseURL),
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

	// Insert using GORM - it will handle JSONB serialization
	h.db.GetDB().Create(&plan)
}

// CreateTestUser creates a user and workspace for testing email already registered scenarios
func (h *OnboardingTestHelper) CreateTestUser(email, password, firstName, lastName string) error {
	workspaceID := "ws_test_" + email
	userID := "user_test_" + email

	// Hash password
	hashedPass, err := hash.Password(password)
	if err != nil {
		return err
	}

	// Create workspace using GORM (workspace table only has id and owner_id)
	workspace := map[string]interface{}{
		"id":       workspaceID,
		"owner_id": userID,
	}
	if err := h.db.GetDB().Table("workspaces").Create(workspace).Error; err != nil {
		return err
	}

	// Create user using GORM
	user := map[string]interface{}{
		"id":                userID,
		"workspace_id":      workspaceID,
		"email":             email,
		"password":          hashedPass,
		"first_name":        firstName,
		"last_name":         lastName,
		"is_email_verified": true,
		"role":              "admin",
	}
	return h.db.GetDB().Table("users").Create(user).Error
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
	hashedOTP, err := hash.Password(otp)
	if err != nil {
		return err
	}

	expiry := time.Now().Add(expiresIn)
	return h.db.GetDB().Exec(`UPDATE onboarding_sessions 
		SET otp_hash = ?, otp_expiry = ?, stage = 'identity_pending' 
		WHERE token = ?`, hashedOTP, expiry, token).Error
}

// SetExpiredOTP sets an expired OTP for testing
func (h *OnboardingTestHelper) SetExpiredOTP(token string) error {
	expiredTime := time.Now().UTC().Add(-20 * time.Minute)
	return h.db.GetDB().Exec(`UPDATE onboarding_sessions SET otp_expiry = ? WHERE token = ?`, expiredTime, token).Error
} // UpdateSessionStage updates session stage
func (h *OnboardingTestHelper) UpdateSessionStage(token, stage string) error {
	return h.db.GetDB().Exec("UPDATE onboarding_sessions SET stage = ? WHERE token = ?", stage, token).Error
}

// ExpireSession sets session to expired
func (h *OnboardingTestHelper) ExpireSession(token string) error {
	return h.db.GetDB().Exec("UPDATE onboarding_sessions SET expires_at = NOW() - INTERVAL '1 hour' WHERE token = ?", token).Error
}

// GetSession retrieves session data
func (h *OnboardingTestHelper) GetSession(token string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	row := h.db.GetDB().Raw("SELECT stage, email, email_verified, otp_expiry FROM onboarding_sessions WHERE token = ?", token).Row()

	var stage, email string
	var emailVerified bool
	var otpExpiry *time.Time
	if err := row.Scan(&stage, &email, &emailVerified, &otpExpiry); err != nil {
		return nil, err
	}

	result["stage"] = stage
	result["email"] = email
	result["emailVerified"] = emailVerified
	result["otpExpiry"] = otpExpiry
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
	token, err := h.CreateOnboardingSession(email, planDescriptor)
	if err != nil || token == "" {
		return "", err
	}

	hashedPass, _ := hash.Password("TestPassword123!")
	err = h.db.GetDB().Exec(`UPDATE onboarding_sessions 
		SET stage = 'identity_verified', 
			email_verified = true,
			first_name = 'Test',
			last_name = 'User',
			password_hash = ?,
			method = 'email'
		WHERE token = ?`, hashedPass, token).Error

	return token, err
} // CreateBusinessStagedSession creates session at ready_to_commit stage
func (h *OnboardingTestHelper) CreateBusinessStagedSession(email, planDescriptor string) (string, error) {
	token, err := h.CreateVerifiedSession(email, planDescriptor)
	if err != nil || token == "" {
		return "", err
	}

	err = h.db.GetDB().Exec(`UPDATE onboarding_sessions 
SET stage = 'ready_to_commit',
business_name = 'Test Business',
business_descriptor = 'test-business',
business_country = 'AE',
business_currency = 'AED',
payment_status = 'skipped'
WHERE token = ?`, token).Error

	return token, err
}
