package account

import (
	"github.com/abdelrahman146/kyora/internal/platform/types/role"
)

// CreateWorkspaceInput represents the request to create a new workspace.
type CreateWorkspaceInput struct {
	OwnerID string `form:"ownerId" json:"ownerId" binding:"required"`
}

// CreateUserInput represents the request to create a new user.
type CreateUserInput struct {
	FirstName string `form:"firstName" json:"firstName" binding:"required"`
	LastName  string `form:"lastName" json:"lastName" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required,email"`
	Password  string `form:"password" json:"password" binding:"required,min=8"`
}

// AcceptInvitationInput represents the request body for accepting an invitation.
type AcceptInvitationInput struct {
	FirstName string `form:"firstName" json:"firstName" binding:"required"`
	LastName  string `form:"lastName" json:"lastName" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required,min=8"`
}

// UpdateUserInput represents the request to update a user.
type UpdateUserInput struct {
	FirstName *string `form:"firstName" json:"firstName"`
	LastName  *string `form:"lastName" json:"lastName"`
}

// InviteUserInput represents the request to invite a user to a workspace.
type InviteUserInput struct {
	Email string    `form:"email" json:"email" binding:"required,email"`
	Role  role.Role `form:"role" json:"role" binding:"required,oneof=user admin"`
}

// UpdateUserRoleInput represents the request to update a user's role.
type UpdateUserRoleInput struct {
	Role role.Role `form:"role" json:"role" binding:"required,oneof=user admin"`
}

// Authentication request types

// loginRequest represents the request to login with email and password.
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// googleLoginRequest represents the request to login with Google OAuth.
type googleLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// refreshRequest represents the request to refresh an access token.
type refreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// forgotPasswordRequest represents the request to initiate password reset.
type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// resetPasswordRequest represents the request to reset password with a token.
type resetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// requestEmailVerificationRequest represents the request to send email verification.
type requestEmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// verifyEmailRequest represents the request to verify email with a token.
type verifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}
