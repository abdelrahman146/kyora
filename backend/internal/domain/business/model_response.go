package business

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/platform/types/asset"
)

// BusinessResponse is the API response for Business entity
type BusinessResponse struct {
	ID                 string                `json:"id"`
	WorkspaceID        string                `json:"workspaceId"`
	Descriptor         string                `json:"descriptor"`
	Name               string                `json:"name"`
	Brand              string                `json:"brand"`
	Logo               *asset.AssetReference `json:"logo,omitempty"`
	CountryCode        string                `json:"countryCode"`
	Currency           string                `json:"currency"`
	StorefrontPublicID string                `json:"storefrontPublicId"`
	StorefrontEnabled  bool                  `json:"storefrontEnabled"`
	StorefrontTheme    StorefrontTheme       `json:"storefrontTheme"`
	SupportEmail       string                `json:"supportEmail"`
	PhoneNumber        string                `json:"phoneNumber"`
	WhatsappNumber     string                `json:"whatsappNumber"`
	Address            string                `json:"address"`
	WebsiteURL         string                `json:"websiteUrl"`
	InstagramURL       string                `json:"instagramUrl"`
	FacebookURL        string                `json:"facebookUrl"`
	TikTokURL          string                `json:"tiktokUrl"`
	XURL               string                `json:"xUrl"`
	SnapchatURL        string                `json:"snapchatUrl"`
	VatRate            string                `json:"vatRate"`
	SafetyBuffer       string                `json:"safetyBuffer"`
	EstablishedAt      time.Time             `json:"establishedAt"`
	ArchivedAt         *time.Time            `json:"archivedAt,omitempty"`
	CreatedAt          time.Time             `json:"createdAt"`
	UpdatedAt          time.Time             `json:"updatedAt"`
}

// ToBusinessResponse converts Business model to BusinessResponse
func ToBusinessResponse(b *Business) BusinessResponse {
	return BusinessResponse{
		ID:                 b.ID,
		WorkspaceID:        b.WorkspaceID,
		Descriptor:         b.Descriptor,
		Name:               b.Name,
		Brand:              b.Brand,
		Logo:               b.Logo,
		CountryCode:        b.CountryCode,
		Currency:           b.Currency,
		StorefrontPublicID: b.StorefrontPublicID,
		StorefrontEnabled:  b.StorefrontEnabled,
		StorefrontTheme:    b.StorefrontTheme,
		SupportEmail:       b.SupportEmail,
		PhoneNumber:        b.PhoneNumber,
		WhatsappNumber:     b.WhatsappNumber,
		Address:            b.Address,
		WebsiteURL:         b.WebsiteURL,
		InstagramURL:       b.InstagramURL,
		FacebookURL:        b.FacebookURL,
		TikTokURL:          b.TikTokURL,
		XURL:               b.XURL,
		SnapchatURL:        b.SnapchatURL,
		VatRate:            b.VatRate.String(),
		SafetyBuffer:       b.SafetyBuffer.StringFixed(2),
		EstablishedAt:      b.EstablishedAt,
		ArchivedAt:         b.ArchivedAt,
		CreatedAt:          b.CreatedAt,
		UpdatedAt:          b.UpdatedAt,
	}
}

// ShippingZoneResponse is the API response for ShippingZone entity
type ShippingZoneResponse struct {
	ID                    string    `json:"id"`
	BusinessID            string    `json:"businessId"`
	Name                  string    `json:"name"`
	Countries             []string  `json:"countries"`
	Currency              string    `json:"currency"`
	ShippingCost          string    `json:"shippingCost"`
	FreeShippingThreshold string    `json:"freeShippingThreshold"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

// ToShippingZoneResponse converts ShippingZone model to ShippingZoneResponse
func ToShippingZoneResponse(z *ShippingZone) ShippingZoneResponse {
	resp := ShippingZoneResponse{
		ID:                    z.ID,
		BusinessID:            z.BusinessID,
		Name:                  z.Name,
		Countries:             []string(z.Countries),
		Currency:              z.Currency,
		ShippingCost:          z.ShippingCost.String(),
		FreeShippingThreshold: z.FreeShippingThreshold.String(),
		CreatedAt:             z.CreatedAt,
		UpdatedAt:             z.UpdatedAt,
	}
	if resp.Countries == nil {
		resp.Countries = []string{}
	}
	return resp
}

// ToShippingZoneResponses converts a slice of ShippingZone models to responses
func ToShippingZoneResponses(zones []ShippingZone) []ShippingZoneResponse {
	responses := make([]ShippingZoneResponse, len(zones))
	for i, zone := range zones {
		responses[i] = ToShippingZoneResponse(&zone)
	}
	return responses
}

// PaymentMethodResponse is the API response for payment methods
type PaymentMethodResponse struct {
	Descriptor        string `json:"descriptor"`
	Name              string `json:"name"`
	LogoURL           string `json:"logoUrl"`
	Enabled           bool   `json:"enabled"`
	FeePercent        string `json:"feePercent"`
	FeeFixed          string `json:"feeFixed"`
	DefaultFeePercent string `json:"defaultFeePercent"`
	DefaultFeeFixed   string `json:"defaultFeeFixed"`
}

// ToPaymentMethodResponse converts BusinessPaymentMethodView to PaymentMethodResponse
func ToPaymentMethodResponse(v BusinessPaymentMethodView) PaymentMethodResponse {
	return PaymentMethodResponse{
		Descriptor:        string(v.Descriptor),
		Name:              v.Name,
		LogoURL:           v.LogoURL,
		Enabled:           v.Enabled,
		FeePercent:        v.FeePercent.String(),
		FeeFixed:          v.FeeFixed.String(),
		DefaultFeePercent: v.DefaultFeePercent.String(),
		DefaultFeeFixed:   v.DefaultFeeFixed.String(),
	}
}

// ToPaymentMethodResponses converts a slice of BusinessPaymentMethodView to responses
func ToPaymentMethodResponses(views []BusinessPaymentMethodView) []PaymentMethodResponse {
	responses := make([]PaymentMethodResponse, len(views))
	for i, view := range views {
		responses[i] = ToPaymentMethodResponse(view)
	}
	return responses
}
