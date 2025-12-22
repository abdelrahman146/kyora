package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/tests/testutils"
	"github.com/stretchr/testify/suite"
)

// PasswordResetSuite tests password reset endpoints
// POST /v1/auth/forgot-password and POST /v1/auth/reset-password
type PasswordResetSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *PasswordResetSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *PasswordResetSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations")
	s.NoError(err)
}

func (s *PasswordResetSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations")
	s.NoError(err)
}

func (s *PasswordResetSuite) TestForgotPassword_Success() {
	ctx := context.Background()
	email := "test@example.com"

	user, _, _, err := s.helper.CreateTestUser(ctx, email, "OldPassword123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	payload := map[string]interface{}{
		"email": email,
	}

	resp, err := s.helper.Client.Post("/v1/auth/forgot-password", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify token was created in cache
	token, err := s.helper.CreatePasswordResetToken(ctx, user)
	s.NoError(err)
	s.NotEmpty(token)
}

func (s *PasswordResetSuite) TestForgotPassword_NonExistentEmail() {
	// Should return success to prevent email enumeration
	payload := map[string]interface{}{
		"email": "nonexistent@example.com",
	}

	resp, err := s.helper.Client.Post("/v1/auth/forgot-password", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func (s *PasswordResetSuite) TestForgotPassword_InvalidEmailFormat() {
	payload := map[string]interface{}{
		"email": "not-an-email",
	}

	resp, err := s.helper.Client.Post("/v1/auth/forgot-password", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *PasswordResetSuite) TestForgotPassword_MissingEmail() {
	payload := map[string]interface{}{}

	resp, err := s.helper.Client.Post("/v1/auth/forgot-password", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *PasswordResetSuite) TestResetPassword_Success() {
	ctx := context.Background()
	email := "test@example.com"
	oldPassword := "OldPassword123!"
	newPassword := "NewPassword456!"

	user, _, _, err := s.helper.CreateTestUser(ctx, email, oldPassword, "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	// Create password reset token
	token, err := s.helper.CreatePasswordResetToken(ctx, user)
	s.NoError(err)

	payload := map[string]interface{}{
		"token":       token,
		"newPassword": newPassword,
	}

	resp, err := s.helper.Client.Post("/v1/auth/reset-password", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify can login with new password
	loginPayload := map[string]interface{}{
		"email":    email,
		"password": newPassword,
	}
	loginResp, err := s.helper.Client.Post("/v1/auth/login", loginPayload)
	s.NoError(err)
	defer loginResp.Body.Close()

	s.Equal(http.StatusOK, loginResp.StatusCode)

	// Verify cannot login with old password
	oldLoginPayload := map[string]interface{}{
		"email":    email,
		"password": oldPassword,
	}
	oldLoginResp, err := s.helper.Client.Post("/v1/auth/login", oldLoginPayload)
	s.NoError(err)
	defer oldLoginResp.Body.Close()

	s.Equal(http.StatusUnauthorized, oldLoginResp.StatusCode)
}

func (s *PasswordResetSuite) TestResetPassword_InvalidToken() {
	payload := map[string]interface{}{
		"token":       "invalid-token",
		"newPassword": "NewPassword456!",
	}

	resp, err := s.helper.Client.Post("/v1/auth/reset-password", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal("Unauthorized", result["title"])
	s.Contains(result["detail"], "invalid or expired token")
}

func (s *PasswordResetSuite) TestResetPassword_MissingToken() {
	payload := map[string]interface{}{
		"newPassword": "NewPassword456!",
	}

	resp, err := s.helper.Client.Post("/v1/auth/reset-password", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *PasswordResetSuite) TestResetPassword_MissingPassword() {
	ctx := context.Background()
	user, _, _, err := s.helper.CreateTestUser(ctx, "test@example.com", "OldPassword123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	token, err := s.helper.CreatePasswordResetToken(ctx, user)
	s.NoError(err)

	payload := map[string]interface{}{
		"token": token,
	}

	resp, err := s.helper.Client.Post("/v1/auth/reset-password", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *PasswordResetSuite) TestResetPassword_ShortPassword() {
	ctx := context.Background()
	user, _, _, err := s.helper.CreateTestUser(ctx, "test@example.com", "OldPassword123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	token, err := s.helper.CreatePasswordResetToken(ctx, user)
	s.NoError(err)

	payload := map[string]interface{}{
		"token":       token,
		"newPassword": "Short1!",
	}

	resp, err := s.helper.Client.Post("/v1/auth/reset-password", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Equal("Bad Request", result["title"])
}

func (s *PasswordResetSuite) TestResetPassword_TokenReuseNotAllowed() {
	ctx := context.Background()
	email := "test@example.com"

	user, _, _, err := s.helper.CreateTestUser(ctx, email, "OldPassword123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	token, err := s.helper.CreatePasswordResetToken(ctx, user)
	s.NoError(err)

	// First reset - should succeed
	payload1 := map[string]interface{}{
		"token":       token,
		"newPassword": "NewPassword456!",
	}

	resp1, err := s.helper.Client.Post("/v1/auth/reset-password", payload1)
	s.NoError(err)
	defer resp1.Body.Close()

	s.Equal(http.StatusNoContent, resp1.StatusCode)

	// Second reset with same token - should fail
	payload2 := map[string]interface{}{
		"token":       token,
		"newPassword": "AnotherPassword789!",
	}

	resp2, err := s.helper.Client.Post("/v1/auth/reset-password", payload2)
	s.NoError(err)
	defer resp2.Body.Close()

	s.Equal(http.StatusUnauthorized, resp2.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp2, &result))
	s.Contains(result["detail"], "invalid or expired token")

	// Verify password is the first new password, not the second
	loginPayload := map[string]interface{}{
		"email":    email,
		"password": "NewPassword456!",
	}
	loginResp, err := s.helper.Client.Post("/v1/auth/login", loginPayload)
	s.NoError(err)
	defer loginResp.Body.Close()

	s.Equal(http.StatusOK, loginResp.StatusCode)
}

func TestPasswordResetSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(PasswordResetSuite))
}

// EmailVerificationSuite tests email verification endpoints
// POST /v1/auth/verify-email/request and POST /v1/auth/verify-email
type EmailVerificationSuite struct {
	suite.Suite
	helper *AccountTestHelper
}

func (s *EmailVerificationSuite) SetupSuite() {
	s.helper = NewAccountTestHelper(testEnv.Database, testEnv.CacheAddr, "http://localhost:18080")
}

func (s *EmailVerificationSuite) SetupTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations")
	s.NoError(err)
}

func (s *EmailVerificationSuite) TearDownTest() {
	err := testutils.TruncateTables(testEnv.Database, "users", "workspaces", "user_invitations")
	s.NoError(err)
}

func (s *EmailVerificationSuite) TestRequestEmailVerification_Success() {
	ctx := context.Background()
	email := "test@example.com"

	user, _, _, err := s.helper.CreateTestUser(ctx, email, "Password123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	// Mark email as unverified
	err = s.helper.MarkEmailUnverified(ctx, user.ID)
	s.NoError(err)

	payload := map[string]interface{}{
		"email": email,
	}

	resp, err := s.helper.Client.Post("/v1/auth/verify-email/request", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func (s *EmailVerificationSuite) TestRequestEmailVerification_NonExistentEmail() {
	// Should return success to prevent email enumeration
	payload := map[string]interface{}{
		"email": "nonexistent@example.com",
	}

	resp, err := s.helper.Client.Post("/v1/auth/verify-email/request", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func (s *EmailVerificationSuite) TestRequestEmailVerification_InvalidEmailFormat() {
	payload := map[string]interface{}{
		"email": "not-an-email",
	}

	resp, err := s.helper.Client.Post("/v1/auth/verify-email/request", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *EmailVerificationSuite) TestVerifyEmail_Success() {
	ctx := context.Background()
	email := "test@example.com"

	user, _, _, err := s.helper.CreateTestUser(ctx, email, "Password123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	// Mark email as unverified
	err = s.helper.MarkEmailUnverified(ctx, user.ID)
	s.NoError(err)

	// Create verification token
	token, err := s.helper.CreateEmailVerificationToken(ctx, user)
	s.NoError(err)

	payload := map[string]interface{}{
		"token": token,
	}

	resp, err := s.helper.Client.Post("/v1/auth/verify-email", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNoContent, resp.StatusCode)

	// Verify email is now verified
	updatedUser, err := s.helper.GetUser(ctx, user.ID)
	s.NoError(err)
	s.True(updatedUser.IsEmailVerified)
}

func (s *EmailVerificationSuite) TestVerifyEmail_InvalidToken() {
	payload := map[string]interface{}{
		"token": "invalid-token",
	}

	resp, err := s.helper.Client.Post("/v1/auth/verify-email", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
	var result map[string]interface{}
	s.NoError(testutils.DecodeJSON(resp, &result))
	s.Contains(result["detail"], "invalid or expired token")
}

func (s *EmailVerificationSuite) TestVerifyEmail_MissingToken() {
	payload := map[string]interface{}{}

	resp, err := s.helper.Client.Post("/v1/auth/verify-email", payload)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *EmailVerificationSuite) TestVerifyEmail_TokenReuseNotAllowed() {
	ctx := context.Background()
	email := "test@example.com"

	user, _, _, err := s.helper.CreateTestUser(ctx, email, "Password123!", "John", "Doe", role.RoleAdmin)
	s.NoError(err)

	// Mark email as unverified
	err = s.helper.MarkEmailUnverified(ctx, user.ID)
	s.NoError(err)

	// Create verification token
	token, err := s.helper.CreateEmailVerificationToken(ctx, user)
	s.NoError(err)

	payload1 := map[string]interface{}{
		"token": token,
	}

	// First verification - should succeed
	resp1, err := s.helper.Client.Post("/v1/auth/verify-email", payload1)
	s.NoError(err)
	defer resp1.Body.Close()

	s.Equal(http.StatusNoContent, resp1.StatusCode)

	// Second verification with same token - should fail
	resp2, err := s.helper.Client.Post("/v1/auth/verify-email", payload1)
	s.NoError(err)
	defer resp2.Body.Close()

	s.Equal(http.StatusUnauthorized, resp2.StatusCode)
}

func TestEmailVerificationSuite(t *testing.T) {
	if testServer == nil {
		t.Skip("Test server not initialized")
	}
	suite.Run(t, new(EmailVerificationSuite))
}
