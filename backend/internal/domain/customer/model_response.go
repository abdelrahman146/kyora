package customer

import (
	"time"
)

// CustomerResponse is the API response for Customer entity
type CustomerResponse struct {
	ID                string         `json:"id"`
	BusinessID        string         `json:"businessId"`
	Name              string         `json:"name"`
	CountryCode       string         `json:"countryCode"`
	Gender            CustomerGender `json:"gender"`
	Email             string         `json:"email"`
	PhoneNumber       string         `json:"phoneNumber,omitempty"`
	PhoneCode         string         `json:"phoneCode,omitempty"`
	TikTokUsername    string         `json:"tiktokUsername,omitempty"`
	InstagramUsername string         `json:"instagramUsername,omitempty"`
	FacebookUsername  string         `json:"facebookUsername,omitempty"`
	XUsername         string         `json:"xUsername,omitempty"`
	SnapchatUsername  string         `json:"snapchatUsername,omitempty"`
	WhatsappNumber    string         `json:"whatsappNumber,omitempty"`
	JoinedAt          time.Time      `json:"joinedAt"`
	OrdersCount       int            `json:"ordersCount"`
	TotalSpent        float64        `json:"totalSpent"`
	AvatarUrl         *string        `json:"avatarUrl,omitempty"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}

// ToCustomerResponse converts Customer model to CustomerResponse
func ToCustomerResponse(c *Customer, ordersCount int, totalSpent float64) CustomerResponse {
	return CustomerResponse{
		ID:                c.ID,
		BusinessID:        c.BusinessID,
		Name:              c.Name,
		CountryCode:       c.CountryCode,
		Gender:            c.Gender,
		Email:             c.Email.String,
		PhoneNumber:       c.PhoneNumber.String,
		PhoneCode:         c.PhoneCode.String,
		TikTokUsername:    c.TikTokUsername.String,
		InstagramUsername: c.InstagramUsername.String,
		FacebookUsername:  c.FacebookUsername.String,
		XUsername:         c.XUsername.String,
		SnapchatUsername:  c.SnapchatUsername.String,
		WhatsappNumber:    c.WhatsappNumber.String,
		JoinedAt:          c.JoinedAt,
		OrdersCount:       ordersCount,
		TotalSpent:        totalSpent,
		AvatarUrl:         nil,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}
}

// ToCustomerResponses converts a slice of Customer models to responses
func ToCustomerResponses(customers []*Customer, ordersCount []int, totalSpent []float64) []CustomerResponse {
	responses := make([]CustomerResponse, len(customers))
	for i, c := range customers {
		oc := 0
		ts := 0.0
		if i < len(ordersCount) {
			oc = ordersCount[i]
		}
		if i < len(totalSpent) {
			ts = totalSpent[i]
		}
		responses[i] = ToCustomerResponse(c, oc, ts)
	}
	return responses
}

// CustomerAddressResponse is the API response for CustomerAddress entity
type CustomerAddressResponse struct {
	ID          string    `json:"id"`
	CustomerID  string    `json:"customerId"`
	CountryCode string    `json:"countryCode"`
	State       string    `json:"state"`
	City        string    `json:"city"`
	Street      string    `json:"street,omitempty"`
	PhoneCode   string    `json:"phoneCode"`
	PhoneNumber string    `json:"phoneNumber"`
	ZipCode     string    `json:"zipCode,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ToCustomerAddressResponse converts CustomerAddress model to CustomerAddressResponse
func ToCustomerAddressResponse(a *CustomerAddress) CustomerAddressResponse {
	return CustomerAddressResponse{
		ID:          a.ID,
		CustomerID:  a.CustomerID,
		CountryCode: a.CountryCode,
		State:       a.State,
		City:        a.City,
		Street:      a.Street.String,
		PhoneCode:   a.PhoneCode,
		PhoneNumber: a.PhoneNumber,
		ZipCode:     a.ZipCode.String,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// ToCustomerAddressResponses converts a slice of CustomerAddress models to responses
func ToCustomerAddressResponses(addresses []*CustomerAddress) []CustomerAddressResponse {
	responses := make([]CustomerAddressResponse, len(addresses))
	for i, a := range addresses {
		responses[i] = ToCustomerAddressResponse(a)
	}
	return responses
}

// CustomerNoteResponse is the API response for CustomerNote entity
type CustomerNoteResponse struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customerId"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ToCustomerNoteResponse converts CustomerNote model to CustomerNoteResponse
func ToCustomerNoteResponse(n *CustomerNote) CustomerNoteResponse {
	return CustomerNoteResponse{
		ID:         n.ID,
		CustomerID: n.CustomerID,
		Content:    n.Content,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
}

// ToCustomerNoteResponses converts a slice of CustomerNote models to responses
func ToCustomerNoteResponses(notes []*CustomerNote) []CustomerNoteResponse {
	responses := make([]CustomerNoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = ToCustomerNoteResponse(n)
	}
	return responses
}
