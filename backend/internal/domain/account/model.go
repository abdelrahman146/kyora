package account

import (
	"database/sql"
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/role"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"gorm.io/gorm"
)

/* Workspace Model */
//-----------------*/

const (
	WorkspaceTable  = "workspaces"
	WorkspaceStruct = "Workspace"
	WorkspacePrefix = "wrk"
)

type Workspace struct {
	gorm.Model
	ID                    string         `gorm:"column:id;primaryKey;type:text" json:"id"`
	OwnerID               string         `gorm:"column:owner_id;type:text" json:"ownerId"`
	StripeCustomerID      sql.NullString `gorm:"column:stripe_customer_id;type:text;unique" json:"stripeCustomerId"`
	StripePaymentMethodID sql.NullString `gorm:"column:stripe_payment_method_id;type:text" json:"stripePaymentMethodId"`
	Users                 []User         `gorm:"foreignKey:WorkspaceID;references:ID" json:"users,omitempty"`
}

func (m *Workspace) TableName() string {
	return WorkspaceTable
}

func (m *Workspace) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(WorkspacePrefix)
	}
	return nil
}

type CreateWorkspaceInput struct {
	OwnerID string `form:"ownerId" json:"ownerId" binding:"required"`
}

var WorkspaceSchema = struct {
	ID                    schema.Field
	OwnerID               schema.Field
	StripeCustomerID      schema.Field
	StripePaymentMethodID schema.Field
	CreatedAt             schema.Field
	UpdatedAt             schema.Field
	DeletedAt             schema.Field
}{
	ID:                    schema.NewField("id", "id"),
	OwnerID:               schema.NewField("owner_id", "ownerId"),
	StripeCustomerID:      schema.NewField("stripe_customer_id", "stripeCustomerId"),
	StripePaymentMethodID: schema.NewField("stripe_payment_method_id", "stripePaymentMethodId"),
	CreatedAt:             schema.NewField("created_at", "createdAt"),
	UpdatedAt:             schema.NewField("updated_at", "updatedAt"),
	DeletedAt:             schema.NewField("deleted_at", "deletedAt"),
}

/* User Model */
//------------*/

const (
	UserTable  = "users"
	UserStruct = "User"
	UserPrefix = "usr"
)

type User struct {
	gorm.Model
	ID              string     `gorm:"column:id;primaryKey;type:text" json:"id"`
	WorkspaceID     string     `gorm:"column:workspace_id;type:text" json:"workspaceId"`
	Workspace       *Workspace `gorm:"foreignKey:WorkspaceID;references:ID" json:"workspace,omitempty"`
	Role            role.Role  `gorm:"column:role;type:text;default:'user'" json:"role"`
	FirstName       string     `gorm:"column:first_name;type:text" json:"firstName"`
	LastName        string     `gorm:"column:last_name;type:text" json:"lastName"`
	Email           string     `gorm:"column:email;type:text;uniqueIndex" json:"email"`
	Password        string     `gorm:"column:password;type:text" json:"-"`
	IsEmailVerified bool       `gorm:"column:is_email_verified;type:boolean;default:false" json:"isEmailVerified"`
	AuthVersion     int        `gorm:"column:auth_version;type:int;default:1" json:"-"`
}

/* Session Model */
//---------------*/

const (
	SessionTable  = "sessions"
	SessionStruct = "Session"
	SessionPrefix = "ses"
)

// Session represents a long-lived server-side authentication session.
//
// A session is created when a user logs in (or completes an auth flow) and is tied to an
// opaque refresh token returned once to the client. Only the refresh token hash is stored.
//
// Security model:
// - Access tokens (JWT) are short-lived and stateless.
// - Sessions are stateful and can be revoked server-side.
// - On sensitive changes (e.g. password reset), we revoke all sessions and bump User.AuthVersion.
type Session struct {
	gorm.Model
	ID          string    `gorm:"column:id;primaryKey;type:text"`
	UserID      string    `gorm:"column:user_id;type:text;index"`
	WorkspaceID string    `gorm:"column:workspace_id;type:text;index"`
	TokenHash   string    `gorm:"column:token_hash;type:text;uniqueIndex"`
	ExpiresAt   time.Time `gorm:"column:expires_at;type:timestamp with time zone;index"`
	CreatedIP   string    `gorm:"column:created_ip;type:text"`
	UserAgent   string    `gorm:"column:user_agent;type:text"`
}

func (m *Session) TableName() string {
	return SessionTable
}

func (m *Session) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(SessionPrefix)
	}
	return nil
}

var SessionSchema = struct {
	ID          schema.Field
	UserID      schema.Field
	WorkspaceID schema.Field
	TokenHash   schema.Field
	ExpiresAt   schema.Field
	CreatedAt   schema.Field
	UpdatedAt   schema.Field
	DeletedAt   schema.Field
}{
	ID:          schema.NewField("id", "id"),
	UserID:      schema.NewField("user_id", "userId"),
	WorkspaceID: schema.NewField("workspace_id", "workspaceId"),
	TokenHash:   schema.NewField("token_hash", "tokenHash"),
	ExpiresAt:   schema.NewField("expires_at", "expiresAt"),
	CreatedAt:   schema.NewField("created_at", "createdAt"),
	UpdatedAt:   schema.NewField("updated_at", "updatedAt"),
	DeletedAt:   schema.NewField("deleted_at", "deletedAt"),
}

func (m *User) TableName() string {
	return UserTable
}

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(UserPrefix)
	}
	if m.AuthVersion <= 0 {
		m.AuthVersion = 1
	}
	return nil
}

type CreateUserInput struct {
	FirstName string `form:"firstName" json:"firstName" binding:"required"`
	LastName  string `form:"lastName" json:"lastName" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required,email"`
	Password  string `form:"password" json:"password" binding:"required,min=8"`
}

// AcceptInvitationInput represents the request body for accepting an invitation
type AcceptInvitationInput struct {
	FirstName string `form:"firstName" json:"firstName" binding:"required"`
	LastName  string `form:"lastName" json:"lastName" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required,min=8"`
}

type UpdateUserInput struct {
	FirstName *string `form:"firstName" json:"firstName"`
	LastName  *string `form:"lastName" json:"lastName"`
}

var UserSchema = struct {
	ID              schema.Field
	WorkspaceID     schema.Field
	Role            schema.Field
	FirstName       schema.Field
	LastName        schema.Field
	Email           schema.Field
	IsEmailVerified schema.Field
	CreatedAt       schema.Field
	UpdatedAt       schema.Field
	DeletedAt       schema.Field
}{
	CreatedAt:       schema.NewField("created_at", "createdAt"),
	UpdatedAt:       schema.NewField("updated_at", "updatedAt"),
	DeletedAt:       schema.NewField("deleted_at", "deletedAt"),
	ID:              schema.NewField("id", "id"),
	WorkspaceID:     schema.NewField("workspace_id", "workspaceId"),
	Role:            schema.NewField("role", "role"),
	FirstName:       schema.NewField("first_name", "firstName"),
	LastName:        schema.NewField("last_name", "lastName"),
	Email:           schema.NewField("email", "email"),
	IsEmailVerified: schema.NewField("is_email_verified", "isEmailVerified"),
}

/* User Invitation Model */
//-----------------------*/

const (
	UserInvitationTable  = "user_invitations"
	UserInvitationStruct = "UserInvitation"
	UserInvitationPrefix = "inv"
)

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusExpired  InvitationStatus = "expired"
	InvitationStatusRevoked  InvitationStatus = "revoked"
)

type UserInvitation struct {
	gorm.Model
	ID          string           `gorm:"column:id;primaryKey;type:text" json:"id"`
	WorkspaceID string           `gorm:"column:workspace_id;type:text" json:"workspaceId"`
	Workspace   *Workspace       `gorm:"foreignKey:WorkspaceID;references:ID" json:"workspace,omitempty"`
	Email       string           `gorm:"column:email;type:text" json:"email"`
	Role        role.Role        `gorm:"column:role;type:text;default:'user'" json:"role"`
	InviterID   string           `gorm:"column:inviter_id;type:text" json:"inviterId"`
	Inviter     *User            `gorm:"foreignKey:InviterID;references:ID" json:"inviter,omitempty"`
	Status      InvitationStatus `gorm:"column:status;type:text;default:'pending'" json:"status"`
	AcceptedAt  *gorm.DeletedAt  `gorm:"column:accepted_at;type:timestamp with time zone" json:"acceptedAt,omitempty"`
}

func (m *UserInvitation) TableName() string {
	return UserInvitationTable
}

func (m *UserInvitation) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(UserInvitationPrefix)
	}
	return nil
}

type InviteUserInput struct {
	Email string    `form:"email" json:"email" binding:"required,email"`
	Role  role.Role `form:"role" json:"role" binding:"required,oneof=user admin"`
}

type UpdateUserRoleInput struct {
	Role role.Role `form:"role" json:"role" binding:"required,oneof=user admin"`
}

var UserInvitationSchema = struct {
	ID          schema.Field
	WorkspaceID schema.Field
	Email       schema.Field
	Role        schema.Field
	InviterID   schema.Field
	Status      schema.Field
	AcceptedAt  schema.Field
	CreatedAt   schema.Field
	UpdatedAt   schema.Field
	DeletedAt   schema.Field
}{
	ID:          schema.NewField("id", "id"),
	WorkspaceID: schema.NewField("workspace_id", "workspaceId"),
	Email:       schema.NewField("email", "email"),
	Role:        schema.NewField("role", "role"),
	InviterID:   schema.NewField("inviter_id", "inviterId"),
	Status:      schema.NewField("status", "status"),
	AcceptedAt:  schema.NewField("accepted_at", "acceptedAt"),
	CreatedAt:   schema.NewField("created_at", "createdAt"),
	UpdatedAt:   schema.NewField("updated_at", "updatedAt"),
	DeletedAt:   schema.NewField("deleted_at", "deletedAt"),
}
