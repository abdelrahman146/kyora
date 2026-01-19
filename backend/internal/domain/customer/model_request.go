package customer

import (
	"time"
)

// CreateCustomerRequest is the request DTO for creating a customer.
type CreateCustomerRequest struct {
	Name              string         `json:"name" binding:"required"`
	Gender            CustomerGender `json:"gender" binding:"omitempty,oneof=male female other"`
	CountryCode       string         `json:"countryCode" binding:"required,len=2"`
	Email             string         `json:"email" binding:"omitempty,email"`
	PhoneNumber       string         `json:"phoneNumber" binding:"required"`
	PhoneCode         string         `json:"phoneCode" binding:"required"`
	TikTokUsername    string         `json:"tiktokUsername" binding:"omitempty"`
	InstagramUsername string         `json:"instagramUsername" binding:"omitempty"`
	FacebookUsername  string         `json:"facebookUsername" binding:"omitempty"`
	XUsername         string         `json:"xUsername" binding:"omitempty"`
	SnapchatUsername  string         `json:"snapchatUsername" binding:"omitempty"`
	WhatsappNumber    string         `json:"whatsappNumber" binding:"omitempty"`
	JoinedAt          time.Time      `json:"joinedAt" binding:"omitempty"`
}

// UpdateCustomerRequest is the request DTO for updating a customer.
type UpdateCustomerRequest struct {
	Name              string         `json:"name" binding:"omitempty"`
	Gender            CustomerGender `json:"gender" binding:"omitempty,oneof=male female other"`
	CountryCode       string         `json:"countryCode" binding:"omitempty,len=2"`
	Email             string         `json:"email" binding:"omitempty,email"`
	PhoneNumber       string         `json:"phoneNumber" binding:"omitempty"`
	PhoneCode         string         `json:"phoneCode" binding:"omitempty"`
	TikTokUsername    string         `json:"tiktokUsername" binding:"omitempty"`
	InstagramUsername string         `json:"instagramUsername" binding:"omitempty"`
	FacebookUsername  string         `json:"facebookUsername" binding:"omitempty"`
	XUsername         string         `json:"xUsername" binding:"omitempty"`
	SnapchatUsername  string         `json:"snapchatUsername" binding:"omitempty"`
	WhatsappNumber    string         `json:"whatsappNumber" binding:"omitempty"`
	JoinedAt          time.Time      `json:"joinedAt" binding:"omitempty"`
}

// CreateCustomerAddressRequest is the request DTO for creating a customer address.
type CreateCustomerAddressRequest struct {
	CountryCode string `json:"countryCode" binding:"required,len=2"`
	State       string `json:"state" binding:"required"`
	City        string `json:"city" binding:"required"`
	PhoneCode   string `json:"phoneCode" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Street      string `json:"street" binding:"omitempty"`
	ZipCode     string `json:"zipCode" binding:"omitempty"`
}

// UpdateCustomerAddressRequest is the request DTO for updating a customer address.
type UpdateCustomerAddressRequest struct {
	Street      string `json:"street" binding:"omitempty"`
	City        string `json:"city" binding:"omitempty"`
	State       string `json:"state" binding:"omitempty"`
	CountryCode string `json:"countryCode" binding:"omitempty,len=2"`
	PhoneCode   string `json:"phoneCode" binding:"omitempty"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty"`
	ZipCode     string `json:"zipCode" binding:"omitempty"`
}

// CreateCustomerNoteRequest is the request DTO for creating a customer note.
type CreateCustomerNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

// UpdateCustomerNoteRequest is the request DTO for updating a customer note.
type UpdateCustomerNoteRequest struct {
	Content string `json:"content" binding:"omitempty"`
}
