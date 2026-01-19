package billing

// attachPMRequest represents a request to attach a payment method to the workspace.
type attachPMRequest struct {
	PaymentMethodID string `json:"paymentMethodId" binding:"required"`
}

// subRequest represents a request to create a subscription.
type subRequest struct {
	PlanDescriptor string `json:"planDescriptor" binding:"required"`
}

// checkoutRequest represents a request to initiate a checkout session.
type checkoutRequest struct {
	PlanDescriptor string `json:"planDescriptor" binding:"required"`
	SuccessURL     string `json:"successUrl" binding:"required,url"`
	CancelURL      string `json:"cancelUrl" binding:"required,url"`
}

// billingPortalRequest represents a request to open the billing portal.
type billingPortalRequest struct {
	ReturnURL string `json:"returnUrl" binding:"required,url"`
}

// scheduleChangeRequest represents a request to schedule a plan change.
type scheduleChangeRequest struct {
	PlanDescriptor string `json:"planDescriptor" binding:"required"`
	EffectiveDate  string `json:"effectiveDate" binding:"required"`  // ISO8601 ("2006-01-02" or "2006-01-02T15:04:05Z")
	ProrationMode  string `json:"prorationMode" binding:"omitempty"` // Stripe proration behavior ("create_prorations" | "none")
}

// prorationEstimateRequest represents a request to estimate proration for a plan change.
type prorationEstimateRequest struct {
	NewPlanDescriptor string `json:"newPlanDescriptor" binding:"required"`
}

// resumeSubscriptionRequest represents a request to resume a subscription.
// Currently no fields; placeholder for future (e.g. payment method enforcement).
type resumeSubscriptionRequest struct {
}

// manualInvoiceRequest represents a request to create a manual invoice.
type manualInvoiceRequest struct {
	Description string  `json:"description" binding:"required"`
	Amount      int64   `json:"amount" binding:"required,min=1"` // amount in minor units (e.g., cents)
	Currency    string  `json:"currency" binding:"required"`
	DueDate     *string `json:"dueDate" binding:"omitempty"` // YYYY-MM-DD
}

// trialExtendRequest represents a request to extend the trial period.
type trialExtendRequest struct {
	AdditionalDays int `json:"additionalDays" binding:"required,min=1,max=30"`
}

// taxCalculateRequest represents a request to calculate tax.
type taxCalculateRequest struct {
	Amount   int64  `json:"amount" binding:"required,min=1"`
	Currency string `json:"currency" binding:"required"`
}
