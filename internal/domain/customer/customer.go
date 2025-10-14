package customer

import (
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"gorm.io/gorm"
)

type CustomerGender string

const (
	GenderMale   CustomerGender = "male"
	GenderFemale CustomerGender = "female"
	GenderOther  CustomerGender = "other"
)

const (
	CustomerTable  = "customers"
	CustomerAlias  = "cus"
	CustomerStruct = "Customer"
)

type Customer struct {
	gorm.Model
	ID                string          `gorm:"column:id;primaryKey;type:text" json:"id"`
	StoreID           string          `gorm:"column:store_id;type:text;not null;index" json:"storeId"`
	Store             *store.Store    `gorm:"foreignKey:StoreID;references:ID" json:"store,omitempty"`
	FirstName         string          `gorm:"column:first_name;type:text;not null" json:"firstName"`
	LastName          string          `gorm:"column:last_name;type:text;not null" json:"lastName"`
	CountryCode       string          `gorm:"column:country_code;type:text;not null" json:"countryCode"`
	Gender            CustomerGender  `gorm:"column:gender;type:text;not null" json:"gender"`
	Email             string          `gorm:"column:email;type:text;not null;uniqueIndex" json:"email"`
	Phone             string          `gorm:"column:phone;type:text" json:"phone,omitempty"`
	TikTokUsername    string          `gorm:"column:tiktok_username;type:text" json:"tiktokUsername,omitempty"`
	InstagramUsername string          `gorm:"column:instagram_username;type:text" json:"instagramUsername,omitempty"`
	FacebookUsername  string          `gorm:"column:facebook_username;type:text" json:"facebookUsername,omitempty"`
	XUsername         string          `gorm:"column:x_username;type:text" json:"xUsername,omitempty"`
	SnapchatUsername  string          `gorm:"column:snapchat_username;type:text" json:"snapchatUsername,omitempty"`
	WhatsappNumber    string          `gorm:"column:whatsapp_number;type:text" json:"whatsappNumber,omitempty"`
	JoinedAt          time.Time       `gorm:"column:joined_at;type:timestamptz;not null;default:now()" json:"joinedAt"`
	Addresses         []*Address      `gorm:"foreignKey:CustomerID;references:ID" json:"addresses,omitempty"`
	Notes             []*CustomerNote `gorm:"foreignKey:CustomerID;references:ID" json:"notes,omitempty"`
}

func (m *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(CustomerAlias)
	}
	return
}

type CreateCustomerRequest struct {
	FirstName         string         `json:"firstName" binding:"required"`
	LastName          string         `json:"lastName"`
	Gender            CustomerGender `json:"gender" binding:"omitempty,oneof=male female other"`
	CountryCode       string         `json:"countryCode" binding:"required,len=2"`
	Email             string         `json:"email" binding:"required,email"`
	Phone             string         `json:"phone" binding:"omitempty"`
	TikTokUsername    string         `json:"tiktokUsername" binding:"omitempty"`
	InstagramUsername string         `json:"instagramUsername" binding:"omitempty"`
	FacebookUsername  string         `json:"facebookUsername" binding:"omitempty"`
	XUsername         string         `json:"xUsername" binding:"omitempty"`
	SnapchatUsername  string         `json:"snapchatUsername" binding:"omitempty"`
	JoinedAt          time.Time      `json:"joinedAt" binding:"omitempty"`
	WhatsappNumber    string         `json:"whatsappNumber" binding:"omitempty"`
}

type UpdateCustomerRequest struct {
	FirstName         string         `json:"firstName" binding:"omitempty"`
	LastName          string         `json:"lastName" binding:"omitempty"`
	Gender            CustomerGender `json:"gender" binding:"omitempty,oneof=male female other"`
	CountryCode       string         `json:"countryCode" binding:"omitempty,len=2"`
	Email             string         `json:"email" binding:"omitempty,email"`
	Phone             string         `json:"phone" binding:"omitempty"`
	TikTokUsername    string         `json:"tiktokUsername" binding:"omitempty"`
	InstagramUsername string         `json:"instagramUsername" binding:"omitempty"`
	FacebookUsername  string         `json:"facebookUsername" binding:"omitempty"`
	XUsername         string         `json:"xUsername" binding:"omitempty"`
	SnapchatUsername  string         `json:"snapchatUsername" binding:"omitempty"`
	WhatsappNumber    string         `json:"whatsappNumber" binding:"omitempty"`
	JoinedAt          time.Time      `json:"joinedAt" binding:"omitempty"`
}
