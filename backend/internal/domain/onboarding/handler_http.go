package onboarding

import (
	"net/http"

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

func (h *HttpHandler) SendEmailOTP(c *gin.Context) {
	var req otpRequest
	if err := request.ValidBody(c, &req); err != nil {
		return
	}
	if err := h.svc.GenerateEmailOTP(c.Request.Context(), req.SessionToken); err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessEmpty(c, http.StatusNoContent)
}

type verifyEmailRequest struct {
	SessionToken string `json:"sessionToken" binding:"required"`
	Code         string `json:"code" binding:"required,len=6"`
	FirstName    string `json:"firstName" binding:"required"`
	LastName     string `json:"lastName" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
}

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
