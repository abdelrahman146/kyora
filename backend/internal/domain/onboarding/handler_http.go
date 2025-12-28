package onboarding

import (
	"net/http"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/request"
	"github.com/abdelrahman146/kyora/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	svc *Service
}

func NewHttpHandler(svc *Service) *HttpHandler {
	h := &HttpHandler{svc: svc}
	return h
}

type startRequest struct {
	Email          string `json:"email" binding:"required,email"`
	PlanDescriptor string `json:"planDescriptor" binding:"required"`
}

type startResponse struct {
	SessionToken string       `json:"sessionToken"`
	Stage        SessionStage `json:"stage"`
	IsPaid       bool         `json:"isPaid"`
}

type stageResponse struct {
	Stage SessionStage `json:"stage"`
}

type paymentStartResponse struct {
	CheckoutURL string `json:"checkoutUrl"`
}

type completeResponse struct {
	User         account.User `json:"user"`
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
}

// Start initializes an onboarding session.
//
// @Summary      Start onboarding session
// @Tags         onboarding
// @Accept       json
// @Produce      json
// @Param        request body startRequest true "Start session request"
// @Success      200 {object} startResponse
// @Failure      400 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/onboarding/start [post]
func (h *HttpHandler) Start(c *gin.Context) {
	var req startRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	sess, err := h.svc.StartSession(c.Request.Context(), req.Email, req.PlanDescriptor, 24*60*60*1e9) // 24h
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"sessionToken": sess.Token, "stage": sess.Stage, "isPaid": sess.IsPaidPlan})
}

type otpRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
}

type otpResponse struct {
	RetryAfterSeconds int `json:"retryAfterSeconds"`
}

// SendEmailOTP generates and sends an OTP to the user's email.
//
// @Summary      Send email OTP
// @Tags         onboarding
// @Accept       json
// @Produce      json
// @Param        request body otpRequest true "OTP request"
// @Success      200 {object} otpResponse
// @Failure      400 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      429 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/onboarding/email/otp [post]
func (h *HttpHandler) SendEmailOTP(c *gin.Context) {
	var req otpRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	retryAfter, err := h.svc.GenerateEmailOTP(c.Request.Context(), req.SessionToken)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, otpResponse{RetryAfterSeconds: int(retryAfter.Seconds())})
}

type verifyEmailRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
	Code         string `json:"code" binding:"required,len=6"`
	FirstName    string `json:"firstName" binding:"required"`
	LastName     string `json:"lastName" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
}

// VerifyEmail verifies an email OTP and stages user profile data.
//
// @Summary      Verify email with OTP
// @Tags         onboarding
// @Accept       json
// @Produce      json
// @Param        request body verifyEmailRequest true "Verify email request"
// @Success      200 {object} stageResponse
// @Failure      400 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/onboarding/email/verify [post]
func (h *HttpHandler) VerifyEmail(c *gin.Context) {
	var req verifyEmailRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	sess, err := h.svc.VerifyEmailOTP(c.Request.Context(), req.SessionToken, req.Code, req.FirstName, req.LastName, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"stage": sess.Stage})
}

type oauthRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
	Code         string `json:"code" binding:"required"`
}

// OAuthGoogle exchanges a Google OAuth code for an identity and stages it on the session.
//
// @Summary      Verify identity with Google OAuth
// @Tags         onboarding
// @Accept       json
// @Produce      json
// @Param        request body oauthRequest true "OAuth Google request"
// @Success      200 {object} stageResponse
// @Failure      400 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/onboarding/oauth/google [post]
func (h *HttpHandler) OAuthGoogle(c *gin.Context) {
	var req oauthRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	sess, err := h.svc.SetOAuthIdentity(c.Request.Context(), req.SessionToken, req.Code)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"stage": sess.Stage})
}

type businessRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Descriptor   string `json:"descriptor" binding:"required"`
	Country      string `json:"country" binding:"required,len=2"`
	Currency     string `json:"currency" binding:"required,len=3"`
}

// SetBusiness stages business details for the session.
//
// @Summary      Set business details
// @Tags         onboarding
// @Accept       json
// @Produce      json
// @Param        request body businessRequest true "Business request"
// @Success      200 {object} stageResponse
// @Failure      400 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      409 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/onboarding/business [post]
func (h *HttpHandler) SetBusiness(c *gin.Context) {
	var req businessRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	sess, err := h.svc.SetBusinessDetails(c.Request.Context(), req.SessionToken, req.Name, req.Descriptor, req.Country, req.Currency)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"stage": sess.Stage})
}

type paymentStartRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
	SuccessURL   string `json:"successUrl" binding:"required,url"`
	CancelURL    string `json:"cancelUrl" binding:"required,url"`
}

// PaymentStart creates a Stripe checkout session for paid tiers.
//
// @Summary      Start payment
// @Tags         onboarding
// @Accept       json
// @Produce      json
// @Param        request body paymentStartRequest true "Payment start request"
// @Success      200 {object} paymentStartResponse
// @Failure      400 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/onboarding/payment/start [post]
func (h *HttpHandler) PaymentStart(c *gin.Context) {
	var req paymentStartRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	url, err := h.svc.InitiatePayment(c.Request.Context(), req.SessionToken, req.SuccessURL, req.CancelURL)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"checkoutUrl": url})
}

type completeRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
}

// Complete finalizes the onboarding and returns JWT tokens.
//
// @Summary      Complete onboarding
// @Tags         onboarding
// @Accept       json
// @Produce      json
// @Param        request body completeRequest true "Complete request"
// @Success      200 {object} completeResponse
// @Failure      400 {object} problem.Problem
// @Failure      401 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/onboarding/complete [post]
func (h *HttpHandler) Complete(c *gin.Context) {
	var req completeRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	user, token, refreshToken, err := h.svc.CompleteOnboarding(c.Request.Context(), req.SessionToken, clientIP, userAgent)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessJSON(c, http.StatusOK, gin.H{"user": user, "token": token, "refreshToken": refreshToken})
}

type sessionResponse struct {
	SessionToken       string         `json:"sessionToken"`
	Email              string         `json:"email"`
	Stage              SessionStage   `json:"stage"`
	EmailVerified      bool           `json:"emailVerified"`
	Method             IdentityMethod `json:"method"`
	FirstName          string         `json:"firstName,omitempty"`
	LastName           string         `json:"lastName,omitempty"`
	PlanID             string         `json:"planId"`
	PlanDescriptor     string         `json:"planDescriptor"`
	IsPaidPlan         bool           `json:"isPaidPlan"`
	BusinessName       string         `json:"businessName,omitempty"`
	BusinessDescriptor string         `json:"businessDescriptor,omitempty"`
	BusinessCountry    string         `json:"businessCountry,omitempty"`
	BusinessCurrency   string         `json:"businessCurrency,omitempty"`
	PaymentStatus      PaymentStatus  `json:"paymentStatus"`
	CheckoutSessionID  string         `json:"checkoutSessionId,omitempty"`
	OTPExpiry          *string        `json:"otpExpiry,omitempty"`
	ExpiresAt          string         `json:"expiresAt"`
}

// GetSession retrieves the current onboarding session state by token.
//
// @Summary      Get onboarding session
// @Tags         onboarding
// @Accept       json
// @Produce      json
// @Param        sessionToken query string true "Session token"
// @Success      200 {object} sessionResponse
// @Failure      400 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/onboarding/session [get]
func (h *HttpHandler) GetSession(c *gin.Context) {
	sessionToken := c.Query("sessionToken")
	if sessionToken == "" {
		response.Error(c, ErrSessionTokenRequired(nil))
		return
	}

	sess, err := h.svc.GetSession(c.Request.Context(), sessionToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	var otpExpiry *string
	if sess.OTPExpiry != nil {
		expStr := sess.OTPExpiry.Format("2006-01-02T15:04:05Z07:00")
		otpExpiry = &expStr
	}

	response.SuccessJSON(c, http.StatusOK, sessionResponse{
		SessionToken:       sess.Token,
		Email:              sess.Email,
		Stage:              sess.Stage,
		EmailVerified:      sess.EmailVerified,
		Method:             sess.Method,
		FirstName:          sess.FirstName,
		LastName:           sess.LastName,
		PlanID:             sess.PlanID,
		PlanDescriptor:     sess.PlanDescriptor,
		IsPaidPlan:         sess.IsPaidPlan,
		BusinessName:       sess.BusinessName,
		BusinessDescriptor: sess.BusinessDescriptor,
		BusinessCountry:    sess.BusinessCountry,
		BusinessCurrency:   sess.BusinessCurrency,
		PaymentStatus:      sess.PaymentStatus,
		CheckoutSessionID:  sess.CheckoutSessionID,
		OTPExpiry:          otpExpiry,
		ExpiresAt:          sess.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type deleteSessionRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
}

// DeleteSession cancels/deletes an onboarding session.
//
// @Summary      Delete onboarding session
// @Tags         onboarding
// @Accept       json
// @Produce      json
// @Param        request body deleteSessionRequest true "Delete session request"
// @Success      204
// @Failure      400 {object} problem.Problem
// @Failure      404 {object} problem.Problem
// @Failure      500 {object} problem.Problem
// @Router       /v1/onboarding/session [delete]
func (h *HttpHandler) DeleteSession(c *gin.Context) {
	var req deleteSessionRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}

	err := h.svc.DeleteSession(c.Request.Context(), req.SessionToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessEmpty(c, http.StatusNoContent)
}
