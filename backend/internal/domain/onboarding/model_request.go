package onboarding

// startRequest represents the request to start an onboarding session.
type startRequest struct {
	Email          string `json:"email" binding:"required,email"`
	PlanDescriptor string `json:"planDescriptor" binding:"required"`
}

// otpRequest represents the request to send an OTP.
type otpRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
}

// verifyEmailRequest represents the request to verify email with OTP.
type verifyEmailRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
	Code         string `json:"code" binding:"required,len=6"`
	FirstName    string `json:"firstName" binding:"required"`
	LastName     string `json:"lastName" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
}

// oauthRequest represents the request to verify identity with OAuth.
type oauthRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
	Code         string `json:"code" binding:"required"`
}

// businessRequest represents the request to set business details.
type businessRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Descriptor   string `json:"descriptor" binding:"required"`
	Country      string `json:"country" binding:"required,len=2"`
	Currency     string `json:"currency" binding:"required,len=3"`
}

// paymentStartRequest represents the request to start payment.
type paymentStartRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
	SuccessURL   string `json:"successUrl" binding:"required,url"`
	CancelURL    string `json:"cancelUrl" binding:"required,url"`
}

// completeRequest represents the request to complete onboarding.
type completeRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
}

// deleteSessionRequest represents the request to delete a session.
type deleteSessionRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
}
