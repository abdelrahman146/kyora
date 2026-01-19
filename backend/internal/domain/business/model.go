package business

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/account"
	"github.com/abdelrahman146/kyora/internal/platform/types/asset"
	"github.com/abdelrahman146/kyora/internal/platform/types/problem"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	BusinessTable    = "businesses"
	BusinessStruct   = "Business"
	BusinessPrefix   = "bus"
	StorefrontPrefix = "sf"
)

// StorefrontTheme is a JSONB-backed theme configuration for the public storefront.
// It is intentionally flexible and optional.
type StorefrontTheme struct {
	PrimaryColor      string `json:"primaryColor"`
	SecondaryColor    string `json:"secondaryColor"`
	AccentColor       string `json:"accentColor"`
	BackgroundColor   string `json:"backgroundColor"`
	TextColor         string `json:"textColor"`
	FontFamily        string `json:"fontFamily"`
	HeadingFontFamily string `json:"headingFontFamily"`
}

func (t StorefrontTheme) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (t *StorefrontTheme) Scan(value any) error {
	if t == nil {
		return problem.InternalError().WithError(errors.New("StorefrontTheme scan into nil receiver"))
	}
	if value == nil {
		*t = StorefrontTheme{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, t)
	case string:
		return json.Unmarshal([]byte(v), t)
	default:
		return problem.InternalError().WithError(errors.New("unexpected scan type for StorefrontTheme"))
	}
}

type Business struct {
	gorm.Model
	ID          string                `gorm:"column:id;primaryKey;type:text" json:"id"`
	Descriptor  string                `gorm:"column:descriptor;type:text;uniqueIndex:idx_workspace_id_descriptor" json:"descriptor"`
	WorkspaceID string                `gorm:"column:workspace_id;type:text;uniqueIndex:idx_workspace_id_descriptor" json:"workspaceId"`
	Workspace   *account.Workspace    `gorm:"foreignKey:WorkspaceID;references:ID" json:"workspace,omitempty"`
	Name        string                `gorm:"column:name;type:text" json:"name"`
	Brand       string                `gorm:"column:brand;type:text" json:"brand"`
	Logo        *asset.AssetReference `gorm:"column:logo;type:jsonb" json:"logo,omitempty"`
	CountryCode string                `gorm:"column:country_code;type:text" json:"countryCode"`
	Currency    string                `gorm:"column:currency;type:text" json:"currency"`

	// Public storefront configuration.
	StorefrontPublicID string          `gorm:"column:storefront_public_id;type:text;uniqueIndex" json:"storefrontPublicId"`
	StorefrontEnabled  bool            `gorm:"column:storefront_enabled;type:boolean;not null;default:false" json:"storefrontEnabled"`
	StorefrontTheme    StorefrontTheme `gorm:"column:storefront_theme;type:jsonb;not null;default:'{}'" json:"storefrontTheme"`

	// Public business details for storefront.
	SupportEmail   string          `gorm:"column:support_email;type:text" json:"supportEmail"`
	PhoneNumber    string          `gorm:"column:phone_number;type:text" json:"phoneNumber"`
	WhatsappNumber string          `gorm:"column:whatsapp_number;type:text" json:"whatsappNumber"`
	Address        string          `gorm:"column:address;type:text" json:"address"`
	WebsiteURL     string          `gorm:"column:website_url;type:text" json:"websiteUrl"`
	InstagramURL   string          `gorm:"column:instagram_url;type:text" json:"instagramUrl"`
	FacebookURL    string          `gorm:"column:facebook_url;type:text" json:"facebookUrl"`
	TikTokURL      string          `gorm:"column:tiktok_url;type:text" json:"tiktokUrl"`
	XURL           string          `gorm:"column:x_url;type:text" json:"xUrl"`
	SnapchatURL    string          `gorm:"column:snapchat_url;type:text" json:"snapchatUrl"`
	VatRate        decimal.Decimal `gorm:"column:vat_rate;type:numeric;not null;default:0" json:"vatRate"`
	SafetyBuffer   decimal.Decimal `gorm:"column:safety_buffer;type:numeric;not null;default:0" json:"safetyBuffer"`
	EstablishedAt  time.Time       `gorm:"column:established_at;type:date;default:now()" json:"establishedAt,omitempty"`
	ArchivedAt     *time.Time      `gorm:"column:archived_at;type:timestamp with time zone" json:"archivedAt,omitempty"`
}

func (m *Business) TableName() string {
	return BusinessTable
}

func (m *Business) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(BusinessPrefix)
	}
	if strings.TrimSpace(m.StorefrontPublicID) == "" {
		m.StorefrontPublicID = id.KsuidWithPrefix(StorefrontPrefix)
	}
	return nil
}

var BusinessSchema = struct {
	ID                 schema.Field
	Descriptor         schema.Field
	WorkspaceID        schema.Field
	Name               schema.Field
	Brand              schema.Field
	Logo               schema.Field
	CountryCode        schema.Field
	Currency           schema.Field
	StorefrontPublicID schema.Field
	StorefrontEnabled  schema.Field
	StorefrontTheme    schema.Field
	SupportEmail       schema.Field
	PhoneNumber        schema.Field
	WhatsappNumber     schema.Field
	Address            schema.Field
	WebsiteURL         schema.Field
	InstagramURL       schema.Field
	FacebookURL        schema.Field
	TikTokURL          schema.Field
	XURL               schema.Field
	SnapchatURL        schema.Field
	VatRate            schema.Field
	SafetyBuffer       schema.Field
	EstablishedAt      schema.Field
	ArchivedAt         schema.Field
	CreatedAt          schema.Field
	UpdatedAt          schema.Field
	DeletedAt          schema.Field
}{
	ID:                 schema.NewField("id", "id"),
	Descriptor:         schema.NewField("descriptor", "descriptor"),
	WorkspaceID:        schema.NewField("workspace_id", "workspaceId"),
	Name:               schema.NewField("name", "name"),
	Brand:              schema.NewField("brand", "brand"),
	Logo:               schema.NewField("logo", "logo"),
	CountryCode:        schema.NewField("country_code", "countryCode"),
	Currency:           schema.NewField("currency", "currency"),
	StorefrontPublicID: schema.NewField("storefront_public_id", "storefrontPublicId"),
	StorefrontEnabled:  schema.NewField("storefront_enabled", "storefrontEnabled"),
	StorefrontTheme:    schema.NewField("storefront_theme", "storefrontTheme"),
	SupportEmail:       schema.NewField("support_email", "supportEmail"),
	PhoneNumber:        schema.NewField("phone_number", "phoneNumber"),
	WhatsappNumber:     schema.NewField("whatsapp_number", "whatsappNumber"),
	Address:            schema.NewField("address", "address"),
	WebsiteURL:         schema.NewField("website_url", "websiteUrl"),
	InstagramURL:       schema.NewField("instagram_url", "instagramUrl"),
	FacebookURL:        schema.NewField("facebook_url", "facebookUrl"),
	TikTokURL:          schema.NewField("tiktok_url", "tiktokUrl"),
	XURL:               schema.NewField("x_url", "xUrl"),
	SnapchatURL:        schema.NewField("snapchat_url", "snapchatUrl"),
	VatRate:            schema.NewField("vat_rate", "vatRate"),
	SafetyBuffer:       schema.NewField("safety_buffer", "safetyBuffer"),
	EstablishedAt:      schema.NewField("established_at", "establishedAt"),
	ArchivedAt:         schema.NewField("archived_at", "archivedAt"),
	CreatedAt:          schema.NewField("created_at", "createdAt"),
	UpdatedAt:          schema.NewField("updated_at", "updatedAt"),
	DeletedAt:          schema.NewField("deleted_at", "deletedAt"),
}

// Business request DTOs are defined in model_request.go

const (
	ShippingZoneTable  = "shipping_zones"
	ShippingZoneStruct = "ShippingZone"
	ShippingZonePrefix = "sz"
)

// CountryCodeList is a JSONB-backed list of ISO 3166-1 alpha-2 country codes.
// Stored as JSONB to keep reads fast and avoid join overhead.
type CountryCodeList []string

func (c CountryCodeList) Value() (driver.Value, error) {
	if c == nil {
		c = []string{}
	}
	b, err := json.Marshal([]string(c))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (c *CountryCodeList) Scan(value any) error {
	if c == nil {
		return problem.InternalError().WithError(errors.New("CountryCodeList scan into nil receiver"))
	}
	if value == nil {
		*c = CountryCodeList{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		var out []string
		if err := json.Unmarshal(v, &out); err != nil {
			return err
		}
		*c = CountryCodeList(out)
		return nil
	case string:
		var out []string
		if err := json.Unmarshal([]byte(v), &out); err != nil {
			return err
		}
		*c = CountryCodeList(out)
		return nil
	default:
		return problem.InternalError().WithError(errors.New("unexpected scan type for CountryCodeList"))
	}
}

func (c CountryCodeList) Contains(countryCode string) bool {
	for _, x := range c {
		if x == countryCode {
			return true
		}
	}
	return false
}

// ShippingZone maps destination countries to shipping pricing rules.
// This is part of the business domain (not a separate domain) to stay DRY/KISS.
type ShippingZone struct {
	gorm.Model
	ID                    string          `gorm:"column:id;primaryKey;type:text" json:"id"`
	BusinessID            string          `gorm:"column:business_id;type:text;not null;index;uniqueIndex:idx_shipping_zone_business_name" json:"businessId"`
	Business              *Business       `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	Name                  string          `gorm:"column:name;type:text;not null;uniqueIndex:idx_shipping_zone_business_name" json:"name"`
	Countries             CountryCodeList `gorm:"column:countries;type:jsonb;not null;default:'[]'" json:"countries"`
	Currency              string          `gorm:"column:currency;type:text;not null" json:"currency"`
	ShippingCost          decimal.Decimal `gorm:"column:shipping_cost;type:numeric;not null;default:0" json:"shippingCost"`
	FreeShippingThreshold decimal.Decimal `gorm:"column:free_shipping_threshold;type:numeric;not null;default:0" json:"freeShippingThreshold"`
}

func (m *ShippingZone) TableName() string { return ShippingZoneTable }

func (m *ShippingZone) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(ShippingZonePrefix)
	}
	return nil
}

var ShippingZoneSchema = struct {
	ID                    schema.Field
	BusinessID            schema.Field
	Name                  schema.Field
	Countries             schema.Field
	Currency              schema.Field
	ShippingCost          schema.Field
	FreeShippingThreshold schema.Field
	CreatedAt             schema.Field
	UpdatedAt             schema.Field
	DeletedAt             schema.Field
}{
	ID:                    schema.NewField("id", "id"),
	BusinessID:            schema.NewField("business_id", "businessId"),
	Name:                  schema.NewField("name", "name"),
	Countries:             schema.NewField("countries", "countries"),
	Currency:              schema.NewField("currency", "currency"),
	ShippingCost:          schema.NewField("shipping_cost", "shippingCost"),
	FreeShippingThreshold: schema.NewField("free_shipping_threshold", "freeShippingThreshold"),
	CreatedAt:             schema.NewField("created_at", "createdAt"),
	UpdatedAt:             schema.NewField("updated_at", "updatedAt"),
	DeletedAt:             schema.NewField("deleted_at", "deletedAt"),
}

// Shipping zone request DTOs are defined in model_request.go
