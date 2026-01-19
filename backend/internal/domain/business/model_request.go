package business

import (
	"github.com/abdelrahman146/kyora/internal/platform/types/asset"
	"github.com/abdelrahman146/kyora/internal/platform/types/date"
	"github.com/shopspring/decimal"
)

// CreateBusinessInput represents the request to create a new business.
type CreateBusinessInput struct {
	Name              string                `form:"name" json:"name" binding:"required"`
	Brand             string                `form:"brand" json:"brand" binding:"omitempty"`
	Logo              *asset.AssetReference `form:"logo" json:"logo" binding:"omitempty"`
	Descriptor        string                `form:"descriptor" json:"descriptor" binding:"required"`
	CountryCode       string                `form:"countryCode" json:"countryCode" binding:"required,len=2"`
	Currency          string                `form:"currency" json:"currency" binding:"required,len=3"`
	StorefrontEnabled bool                  `form:"storefrontEnabled" json:"storefrontEnabled" binding:"omitempty"`
	StorefrontTheme   StorefrontTheme       `form:"storefrontTheme" json:"storefrontTheme" binding:"omitempty"`
	SupportEmail      string                `form:"supportEmail" json:"supportEmail" binding:"omitempty,email"`
	PhoneNumber       string                `form:"phoneNumber" json:"phoneNumber" binding:"omitempty"`
	WhatsappNumber    string                `form:"whatsappNumber" json:"whatsappNumber" binding:"omitempty"`
	Address           string                `form:"address" json:"address" binding:"omitempty"`
	WebsiteURL        string                `form:"websiteUrl" json:"websiteUrl" binding:"omitempty,url"`
	InstagramURL      string                `form:"instagramUrl" json:"instagramUrl" binding:"omitempty,url"`
	FacebookURL       string                `form:"facebookUrl" json:"facebookUrl" binding:"omitempty,url"`
	TikTokURL         string                `form:"tiktokUrl" json:"tiktokUrl" binding:"omitempty,url"`
	XURL              string                `form:"xUrl" json:"xUrl" binding:"omitempty,url"`
	SnapchatURL       string                `form:"snapchatUrl" json:"snapchatUrl" binding:"omitempty,url"`
	VatRate           decimal.Decimal       `form:"vatRate" json:"vatRate" binding:"omitempty"`
	SafetyBuffer      decimal.Decimal       `form:"safetyBuffer" json:"safetyBuffer" binding:"omitempty"`
	EstablishedAt     date.Date             `form:"establishedAt" json:"establishedAt" binding:"omitempty"`
}

// UpdateBusinessInput represents the request to update a business.
type UpdateBusinessInput struct {
	Name              *string               `form:"name" json:"name" binding:"omitempty"`
	Brand             *string               `form:"brand" json:"brand" binding:"omitempty"`
	Logo              *asset.AssetReference `form:"logo" json:"logo" binding:"omitempty"`
	Descriptor        *string               `form:"descriptor" json:"descriptor" binding:"omitempty"`
	CountryCode       *string               `form:"countryCode" json:"countryCode" binding:"omitempty,len=2"`
	Currency          *string               `form:"currency" json:"currency" binding:"omitempty,len=3"`
	StorefrontEnabled *bool                 `form:"storefrontEnabled" json:"storefrontEnabled" binding:"omitempty"`
	StorefrontTheme   *StorefrontTheme      `form:"storefrontTheme" json:"storefrontTheme" binding:"omitempty"`
	SupportEmail      *string               `form:"supportEmail" json:"supportEmail" binding:"omitempty,email"`
	PhoneNumber       *string               `form:"phoneNumber" json:"phoneNumber" binding:"omitempty"`
	WhatsappNumber    *string               `form:"whatsappNumber" json:"whatsappNumber" binding:"omitempty"`
	Address           *string               `form:"address" json:"address" binding:"omitempty"`
	WebsiteURL        *string               `form:"websiteUrl" json:"websiteUrl" binding:"omitempty,url"`
	InstagramURL      *string               `form:"instagramUrl" json:"instagramUrl" binding:"omitempty,url"`
	FacebookURL       *string               `form:"facebookUrl" json:"facebookUrl" binding:"omitempty,url"`
	TikTokURL         *string               `form:"tiktokUrl" json:"tiktokUrl" binding:"omitempty,url"`
	XURL              *string               `form:"xUrl" json:"xUrl" binding:"omitempty,url"`
	SnapchatURL       *string               `form:"snapchatUrl" json:"snapchatUrl" binding:"omitempty,url"`
	VatRate           decimal.NullDecimal   `form:"vatRate" json:"vatRate" binding:"omitempty"`
	SafetyBuffer      decimal.NullDecimal   `form:"safetyBuffer" json:"safetyBuffer" binding:"omitempty"`
	EstablishedAt     *date.Date            `form:"establishedAt" json:"establishedAt" binding:"omitempty"`
}

// CreateShippingZoneRequest represents the request to create a shipping zone.
type CreateShippingZoneRequest struct {
	Name                  string          `json:"name" binding:"required"`
	Countries             []string        `json:"countries" binding:"required,min=1,max=50,dive,required,len=2"`
	ShippingCost          decimal.Decimal `json:"shippingCost" binding:"omitempty"`
	FreeShippingThreshold decimal.Decimal `json:"freeShippingThreshold" binding:"omitempty"`
}

// UpdateShippingZoneRequest represents the request to update a shipping zone.
type UpdateShippingZoneRequest struct {
	Name                  *string             `json:"name" binding:"omitempty"`
	Countries             []string            `json:"countries" binding:"omitempty,min=1,max=50,dive,required,len=2"`
	ShippingCost          decimal.NullDecimal `json:"shippingCost" binding:"omitempty"`
	FreeShippingThreshold decimal.NullDecimal `json:"freeShippingThreshold" binding:"omitempty"`
}
