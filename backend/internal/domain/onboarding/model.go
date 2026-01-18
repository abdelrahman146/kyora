package onboarding

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"gorm.io/gorm"
)

const (
	SessionTable  = "onboarding_sessions"
	SessionStruct = "OnboardingSession"
	SessionPrefix = "obs"
)

// Flow stages – we deliberately keep them linear but allow re-entry before commit
type SessionStage string

const (
	StagePlanSelected     SessionStage = "plan_selected"
	StageIdentityPending  SessionStage = "identity_pending" // after plan, before verification
	StageIdentityVerified SessionStage = "identity_verified"
	StageBusinessStaged   SessionStage = "business_staged"
	StagePaymentPending   SessionStage = "payment_pending" // only for paid tiers
	StageReadyToCommit    SessionStage = "ready_to_commit" // all prerequisites satisfied
	StageCommitted        SessionStage = "committed"       // finalized & removed logically
)

// PaymentStatus tracks paid tier progress; free plans skip to succeeded immediately
type PaymentStatus string

const (
	PaymentStatusSkipped   PaymentStatus = "skipped"   // free plan
	PaymentStatusPending   PaymentStatus = "pending"   // checkout created, waiting
	PaymentStatusSucceeded PaymentStatus = "succeeded" // confirmed via webhook / synchronous return
)

// Method of identity establishment
type IdentityMethod string

const (
	IdentityEmail  IdentityMethod = "email"
	IdentityGoogle IdentityMethod = "google"
)

// OnboardingSession stores staged data; NOTHING is considered permanent until committed.
// All personally identifiable info is kept here until commit then row is deleted.
type OnboardingSession struct {
	gorm.Model
	ID                 string         `gorm:"column:id;primaryKey;type:text" json:"id"`
	Token              string         `gorm:"column:token;type:text;uniqueIndex" json:"token"` // public facing reference
	Email              string         `gorm:"column:email;type:text;uniqueIndex" json:"email"`
	EmailVerified      bool           `gorm:"column:email_verified;type:boolean;default:false" json:"emailVerified"`
	Method             IdentityMethod `gorm:"column:method;type:text" json:"method"`
	FirstName          string         `gorm:"column:first_name;type:text" json:"firstName"`
	LastName           string         `gorm:"column:last_name;type:text" json:"lastName"`
	PasswordHash       string         `gorm:"column:password_hash;type:text" json:"-"`
	PlanID             string         `gorm:"column:plan_id;type:text" json:"planId"`
	PlanDescriptor     string         `gorm:"column:plan_descriptor;type:text" json:"planDescriptor"`
	IsPaidPlan         bool           `gorm:"column:is_paid_plan;type:boolean" json:"isPaidPlan"`
	BusinessName       string         `gorm:"column:business_name;type:text" json:"businessName"`
	BusinessDescriptor string         `gorm:"column:business_descriptor;type:text" json:"businessDescriptor"`
	BusinessCountry    string         `gorm:"column:business_country;type:text" json:"businessCountry"`
	BusinessCurrency   string         `gorm:"column:business_currency;type:text" json:"businessCurrency"`
	StripeCustomerID   string         `gorm:"column:stripe_customer_id;type:text" json:"stripeCustomerId"`
	StripeSubID        string         `gorm:"column:stripe_sub_id;type:text" json:"stripeSubId"`
	CheckoutSessionID  string         `gorm:"column:checkout_session_id;type:text" json:"checkoutSessionId"`
	PaymentStatus      PaymentStatus  `gorm:"column:payment_status;type:text" json:"paymentStatus"`
	Stage              SessionStage   `gorm:"column:stage;type:text;index" json:"stage"`
	OTPHash            string         `gorm:"column:otp_hash;type:text" json:"-"`
	OTPExpiry          *time.Time     `gorm:"column:otp_expiry;type:timestamp" json:"otpExpiry,omitempty"`
	ExpiresAt          time.Time      `gorm:"column:expires_at;type:timestamp;index" json:"expiresAt"`
	CommittedAt        *time.Time     `gorm:"column:committed_at;type:timestamp" json:"committedAt,omitempty"`
}

func (m *OnboardingSession) TableName() string { return SessionTable }

func (m *OnboardingSession) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(SessionPrefix)
	}
	if m.Token == "" {
		if t, err := id.RandomString(40); err == nil {
			m.Token = t
		}
	}
	return nil
}

// Simple schema for ordering/filtering (subset of fields)
var SessionSchema = struct {
	ID            schema.Field
	Token         schema.Field
	Email         schema.Field
	Stage         schema.Field
	CreatedAt     schema.Field
	ExpiresAt     schema.Field
	CommittedAt   schema.Field
	PaymentStatus schema.Field
}{
	ID:            schema.NewField("id", "id"),
	Token:         schema.NewField("token", "token"),
	Email:         schema.NewField("email", "email"),
	Stage:         schema.NewField("stage", "stage"),
	CreatedAt:     schema.NewField("created_at", "createdAt"),
	ExpiresAt:     schema.NewField("expires_at", "expiresAt"),
	CommittedAt:   schema.NewField("committed_at", "committedAt"),
	PaymentStatus: schema.NewField("payment_status", "paymentStatus"),
}
