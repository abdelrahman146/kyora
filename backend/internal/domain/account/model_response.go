package account

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
)

/* User Response DTO */
//-------------------*/

// UserResponse represents the API response shape for a User.
// It excludes GORM metadata and sensitive fields (Password, AuthVersion).
type UserResponse struct {
	ID              string    `json:"id"`
	WorkspaceID     string    `json:"workspaceId"`
	Role            role.Role `json:"role"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// ToUserResponse converts a User model to a UserResponse DTO
func ToUserResponse(user *User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:              user.ID,
		WorkspaceID:     user.WorkspaceID,
		Role:            user.Role,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

// ToUserResponses converts a slice of User models to UserResponse DTOs
func ToUserResponses(users []*User) []*UserResponse {
	if users == nil {
		return nil
	}
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}

/* Workspace Response DTO */
//------------------------*/

// WorkspaceResponse represents the API response shape for a Workspace.
// It excludes GORM metadata and includes only essential fields.
type WorkspaceResponse struct {
	ID                    string    `json:"id"`
	OwnerID               string    `json:"ownerId"`
	StripeCustomerID      *string   `json:"stripeCustomerId,omitempty"`
	StripePaymentMethodID *string   `json:"stripePaymentMethodId,omitempty"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// ToWorkspaceResponse converts a Workspace model to a WorkspaceResponse DTO
func ToWorkspaceResponse(workspace *Workspace) *WorkspaceResponse {
	if workspace == nil {
		return nil
	}

	resp := &WorkspaceResponse{
		ID:        workspace.ID,
		OwnerID:   workspace.OwnerID,
		CreatedAt: workspace.CreatedAt,
		UpdatedAt: workspace.UpdatedAt,
	}

	if workspace.StripeCustomerID.Valid {
		resp.StripeCustomerID = &workspace.StripeCustomerID.String
	}
	if workspace.StripePaymentMethodID.Valid {
		resp.StripePaymentMethodID = &workspace.StripePaymentMethodID.String
	}

	return resp
}

/* User Invitation Response DTO */
//-------------------------------*/

// UserInvitationResponse represents the API response shape for a UserInvitation.
// It excludes GORM metadata and only includes relevant fields.
type UserInvitationResponse struct {
	ID          string           `json:"id"`
	WorkspaceID string           `json:"workspaceId"`
	Email       string           `json:"email"`
	Role        role.Role        `json:"role"`
	InviterID   string           `json:"inviterId"`
	Status      InvitationStatus `json:"status"`
	AcceptedAt  *time.Time       `json:"acceptedAt,omitempty"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

// ToUserInvitationResponse converts a UserInvitation model to a UserInvitationResponse DTO
func ToUserInvitationResponse(invitation *UserInvitation) *UserInvitationResponse {
	if invitation == nil {
		return nil
	}

	resp := &UserInvitationResponse{
		ID:          invitation.ID,
		WorkspaceID: invitation.WorkspaceID,
		Email:       invitation.Email,
		Role:        invitation.Role,
		InviterID:   invitation.InviterID,
		Status:      invitation.Status,
		CreatedAt:   invitation.CreatedAt,
		UpdatedAt:   invitation.UpdatedAt,
	}

	if invitation.AcceptedAt != nil && invitation.AcceptedAt.Valid {
		acceptedTime := invitation.AcceptedAt.Time
		resp.AcceptedAt = &acceptedTime
	}

	return resp
}

// ToUserInvitationResponses converts a slice of UserInvitation models to UserInvitationResponse DTOs
func ToUserInvitationResponses(invitations []*UserInvitation) []*UserInvitationResponse {
	if invitations == nil {
		return nil
	}
	responses := make([]*UserInvitationResponse, len(invitations))
	for i, invitation := range invitations {
		responses[i] = ToUserInvitationResponse(invitation)
	}
	return responses
}

/* Auth Response Types */
//-----------------------*/

// LoginResponse represents the API response shape for login operations.
type LoginResponse struct {
	User         *UserResponse `json:"user"`
	Token        string        `json:"token"`
	RefreshToken string        `json:"refreshToken"`
}

// RefreshResponse represents the API response shape for token refresh operations.
type RefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// ToLoginResponse converts a User model and tokens to a LoginResponse
func ToLoginResponse(user *User, token, refreshToken string) *LoginResponse {
	return &LoginResponse{
		User:         ToUserResponse(user),
		Token:        token,
		RefreshToken: refreshToken,
	}
}
