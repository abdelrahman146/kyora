package testutils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/domain/billing"
	"github.com/abdelrahman146/kyora/internal/platform/auth"
	"github.com/abdelrahman146/kyora/internal/platform/database"
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/utils/hash"
	"github.com/shopspring/decimal"
)

// CreateAuthenticatedUser creates a workspace with a user and returns user, workspace, and JWT token
func CreateAuthenticatedUser(ctx context.Context, db *database.Database, email, password, firstName, lastName string, userRole role.Role) (*account.User, *account.Workspace, string, error) {
	// Hash password
	hashedPassword, err := hash.Password(password)
	if err != nil {
		return nil, nil, "", err
	}

	// Create workspace
	workspace := &account.Workspace{}
	workspaceRepo := database.NewRepository[account.Workspace](db)
	if err := workspaceRepo.CreateOne(ctx, workspace); err != nil {
		return nil, nil, "", err
	}

	// Create user
	user := &account.User{
		WorkspaceID:     workspace.ID,
		Role:            userRole,
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		Password:        hashedPassword,
		IsEmailVerified: true,
	}
	userRepo := database.NewRepository[account.User](db)
	if err := userRepo.CreateOne(ctx, user); err != nil {
		return nil, nil, "", err
	}

	// Set workspace owner
	workspace.OwnerID = user.ID
	if err := workspaceRepo.UpdateOne(ctx, workspace); err != nil {
		return nil, nil, "", err
	}

	// Generate JWT token
	token, err := auth.NewJwtToken(user.ID, user.WorkspaceID, user.AuthVersion)
	if err != nil {
		return nil, nil, "", err
	}

	return user, workspace, token, nil
}

// CreateWorkspaceWithUsers creates a workspace with multiple users
func CreateWorkspaceWithUsers(ctx context.Context, db *database.Database, ownerEmail, ownerPassword string, additionalUsers []struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Role      role.Role
}) (*account.Workspace, []*account.User, error) {
	// Create owner
	owner, workspace, _, err := CreateAuthenticatedUser(ctx, db, ownerEmail, ownerPassword, "Owner", "User", role.RoleAdmin)
	if err != nil {
		return nil, nil, err
	}

	users := []*account.User{owner}

	// Create additional users
	userRepo := database.NewRepository[account.User](db)
	for _, userData := range additionalUsers {
		hashedPassword, err := hash.Password(userData.Password)
		if err != nil {
			return nil, nil, err
		}

		user := &account.User{
			WorkspaceID:     workspace.ID,
			Role:            userData.Role,
			FirstName:       userData.FirstName,
			LastName:        userData.LastName,
			Email:           userData.Email,
			Password:        hashedPassword,
			IsEmailVerified: true,
		}

		if err := userRepo.CreateOne(ctx, user); err != nil {
			return nil, nil, err
		}
		users = append(users, user)
	}

	return workspace, users, nil
}

// CreatePasswordResetToken creates a password reset token for a user
func CreatePasswordResetToken(ctx context.Context, storage *account.Storage, user *account.User) (string, error) {
	payload := &account.PasswordResetPayload{
		UserID:      user.ID,
		WorkspaceID: user.WorkspaceID,
		Email:       user.Email,
	}
	token, _, err := storage.CreatePasswordResetToken(payload)
	return token, err
}

// CreateEmailVerificationToken creates an email verification token for a user
func CreateEmailVerificationToken(ctx context.Context, storage *account.Storage, user *account.User) (string, error) {
	payload := &account.VerifyEmailPayload{
		UserID:      user.ID,
		WorkspaceID: user.WorkspaceID,
		Email:       user.Email,
	}
	token, _, err := storage.CreateVerifyEmailToken(payload)
	return token, err
}

// CreateInvitationToken creates a workspace invitation token
func CreateInvitationToken(ctx context.Context, storage *account.Storage, invitation *account.UserInvitation, inviterID string) (string, error) {
	payload := &account.WorkspaceInvitationPayload{
		InvitationID: invitation.ID,
		WorkspaceID:  invitation.WorkspaceID,
		Email:        invitation.Email,
		Role:         string(invitation.Role),
		InviterID:    inviterID,
	}
	token, _, err := storage.CreateWorkspaceInvitationToken(payload)
	return token, err
}

// LoginAndGetToken makes a login request and returns the JWT token
func LoginAndGetToken(client *HTTPClient, email, password string) (string, error) {
	payload := map[string]interface{}{
		"email":    email,
		"password": password,
	}
	resp, err := client.Post("/v1/auth/login", payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed: status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := DecodeJSON(resp, &result); err != nil {
		return "", err
	}

	if token, ok := result["token"].(string); ok && token != "" {
		return token, nil
	}
	return "", fmt.Errorf("login response missing token")
}

// CreateTestPlan creates a complete test plan with all required fields
func CreateTestPlan(ctx context.Context, db *database.Database, descriptor string) (*billing.Plan, error) {
	planRepo := database.NewRepository[billing.Plan](db)

	// Check if plan already exists
	existingPlan, err := planRepo.FindOne(ctx, planRepo.ScopeEquals(billing.PlanSchema.Descriptor, descriptor))
	if err == nil && existingPlan != nil {
		return existingPlan, nil
	}

	// Create features JSONB
	features := billing.PlanFeature{
		CustomerManagement:       true,
		InventoryManagement:      true,
		OrderManagement:          true,
		ExpenseManagement:        true,
		Accounting:               true,
		BasicAnalytics:           true,
		FinancialReports:         true,
		DataImport:               false,
		DataExport:               false,
		AdvancedAnalytics:        false,
		AdvancedFinancialReports: false,
		OrderPaymentLinks:        false,
		InvoiceGeneration:        false,
		ExportAnalyticsData:      false,
		AIBusinessAssistant:      false,
	}

	// Create limits JSONB
	// Use reasonable limits for testing - allow multiple team members and adequate orders
	limits := billing.PlanLimit{
		MaxOrdersPerMonth: 1000, // High limit for testing
		MaxTeamMembers:    10,   // Allow testing of team features
		MaxBusinesses:     5,    // Allow multiple businesses for testing
	}

	// Generate unique StripePlanID to avoid unique constraint violations
	stripePlanID := "stripe_test_plan_" + descriptor + "_" + time.Now().Format("20060102150405")

	// Create plan
	plan := &billing.Plan{
		Descriptor:   descriptor,
		Name:         "Test " + descriptor + " Plan",
		Description:  "Test plan for E2E testing with descriptor " + descriptor,
		StripePlanID: stripePlanID,
		Price:        decimal.NewFromInt(0),
		Currency:     "aed",
		BillingCycle: billing.BillingCycleMonthly,
		Features:     features,
		Limits:       limits,
	}

	if err := planRepo.CreateOne(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// CreateTestSubscription creates a subscription for a workspace with a test plan
func CreateTestSubscription(ctx context.Context, db *database.Database, workspaceID string) error {
	// Create a test plan first
	plan, err := CreateTestPlan(ctx, db, "test_starter")
	if err != nil {
		return err
	}

	// Create subscription using repository
	subscriptionRepo := database.NewRepository[billing.Subscription](db)

	// Check if subscription already exists for this workspace
	existingSub, err := subscriptionRepo.FindOne(ctx, subscriptionRepo.ScopeWorkspaceID(workspaceID))
	if err != nil && !database.IsRecordNotFound(err) {
		return fmt.Errorf("failed to check existing subscription: %w", err)
	}
	if err == nil && existingSub != nil && existingSub.ID != "" {
		return nil // Subscription already exists
	}

	subscription := &billing.Subscription{
		WorkspaceID:      workspaceID,
		PlanID:           plan.ID,
		Plan:             plan, // Set the plan relationship
		StripeSubID:      "stripe_sub_test_" + workspaceID,
		CurrentPeriodEnd: time.Now().UTC().Add(30 * 24 * time.Hour),
		Status:           billing.SubscriptionStatusActive,
	}

	err = subscriptionRepo.CreateOne(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}
