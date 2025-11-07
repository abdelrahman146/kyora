package email

// Parameter structs for different email notification types

// ForgotPasswordParams contains data for forgot password email
type ForgotPasswordParams struct {
	Email      string
	UserName   string
	ResetURL   string
	ExpiryTime string // e.g., "24 hours"
}

// PasswordResetConfirmationParams contains data for password reset confirmation email
type PasswordResetConfirmationParams struct {
	Email         string
	UserName      string
	LoginURL      string
	ResetDate     string
	ResetTime     string
	ResetLocation string
	ResetIP       string
}

// EmailVerificationParams contains data for email verification email
type EmailVerificationParams struct {
	Email      string
	UserName   string
	VerifyURL  string
	ExpiryTime string // e.g., "24 hours"
}

// WelcomeParams contains data for welcome email
type WelcomeParams struct {
	Email        string
	UserName     string
	DashboardURL string
	GuideURL     string
}

// SubscriptionWelcomeParams contains data for subscription welcome email
type SubscriptionWelcomeParams struct {
	Email           string
	UserName        string
	PlanName        string
	Amount          string
	BillingCycle    string
	NextBillingDate string
	LastFour        string
	DashboardURL    string
	BillingURL      string
}

// PaymentFailedParams contains data for payment failed email
type PaymentFailedParams struct {
	Email            string
	UserName         string
	Amount           string
	PlanName         string
	LastFour         string
	AttemptDate      string
	NextAttemptDate  string
	GracePeriod      string
	UpdatePaymentURL string
	RetryPaymentURL  string
}

// SubscriptionCanceledParams contains data for subscription canceled email
type SubscriptionCanceledParams struct {
	Email           string
	UserName        string
	PlanName        string
	CancelDate      string
	AccessUntilDate string
	RefundAmount    string
	ReactivateURL   string
	DashboardURL    string
	FeedbackURL     string
}

// TrialEndingParams contains data for trial ending email
type TrialEndingParams struct {
	Email           string
	UserName        string
	PlanName        string
	TrialPeriod     string
	DaysRemaining   int
	TrialEndDate    string
	TrialStartDate  string
	FeaturesUsed    string
	ProjectsCreated string
	MonthlyPrice    string
	SubscribeURL    string
	PlansURL        string
}

// PaymentSucceededParams contains data for payment succeeded email (for future use)
type PaymentSucceededParams struct {
	Email           string
	UserName        string
	Amount          string
	PlanName        string
	PaymentDate     string
	InvoiceNumber   string
	NextBillingDate string
	InvoiceURL      string
	DashboardURL    string
}

// SubscriptionUpdatedParams contains data for subscription updated email (for future use)
type SubscriptionUpdatedParams struct {
	Email           string
	UserName        string
	OldPlanName     string
	NewPlanName     string
	EffectiveDate   string
	ProrationAmount string
	NextBillingDate string
	DashboardURL    string
	BillingURL      string
}

// InvoiceGeneratedParams contains data for invoice generated email (for future use)
type InvoiceGeneratedParams struct {
	Email         string
	UserName      string
	InvoiceNumber string
	Amount        string
	DueDate       string
	InvoiceURL    string
	PaymentURL    string
	DashboardURL  string
}

// LoginNotificationParams contains data for login notification email
type LoginNotificationParams struct {
	Email         string
	UserName      string
	LoginDate     string
	LoginTime     string
	LoginLocation string
	LoginIP       string
	DeviceInfo    string
	ResetURL      string
	SupportEmail  string
}

// SubscriptionConfirmedParams contains data for subscription confirmation email (first payment confirmed)
type SubscriptionConfirmedParams struct {
	Email           string
	UserName        string
	PlanName        string
	Amount          string
	PaymentDate     string
	InvoiceNumber   string
	InvoiceURL      string
	NextBillingDate string
	DashboardURL    string
	BillingURL      string
	SupportEmail    string
}
