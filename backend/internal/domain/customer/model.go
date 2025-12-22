package customer

import (
	"database/sql"
	"time"

	"github.com/abdelrahman146/kyora/internal/domain/business"
	"github.com/abdelrahman146/kyora/internal/platform/types/schema"
	"github.com/abdelrahman146/kyora/internal/platform/utils/id"
	"gorm.io/gorm"
)

/* customer model */
//-------------------*/

type CustomerGender string

const (
	GenderMale   CustomerGender = "male"
	GenderFemale CustomerGender = "female"
	GenderOther  CustomerGender = "other"
)

const (
	CustomerTable  = "customers"
	CustomerStruct = "Customer"
	CustomerPrefix = "cus"
)

type Customer struct {
	gorm.Model
	ID                string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	BusinessID        string             `gorm:"column:business_id;type:text;not null;index" json:"businessId"`
	Business          *business.Business `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	Name              string             `gorm:"column:name;type:text;not null" json:"name"`
	CountryCode       string             `gorm:"column:country_code;type:text;not null" json:"countryCode"`
	Gender            CustomerGender     `gorm:"column:gender;type:text;not null" json:"gender"`
	Email             sql.NullString     `gorm:"column:email;type:text;not null;uniqueIndex" json:"email"`
	PhoneNumber       sql.NullString     `gorm:"column:phone_number;type:text" json:"phoneNumber,omitempty"`
	PhoneCode         sql.NullString     `gorm:"column:phone_code;type:text" json:"phoneCode,omitempty"`
	TikTokUsername    sql.NullString     `gorm:"column:tiktok_username;type:text" json:"tiktokUsername,omitempty"`
	InstagramUsername sql.NullString     `gorm:"column:instagram_username;type:text" json:"instagramUsername,omitempty"`
	FacebookUsername  sql.NullString     `gorm:"column:facebook_username;type:text" json:"facebookUsername,omitempty"`
	XUsername         sql.NullString     `gorm:"column:x_username;type:text" json:"xUsername,omitempty"`
	SnapchatUsername  sql.NullString     `gorm:"column:snapchat_username;type:text" json:"snapchatUsername,omitempty"`
	WhatsappNumber    sql.NullString     `gorm:"column:whatsapp_number;type:text" json:"whatsappNumber,omitempty"`
	JoinedAt          time.Time          `gorm:"column:joined_at;type:timestamptz;not null;default:now()" json:"joinedAt"`
	Addresses         []*CustomerAddress `gorm:"foreignKey:CustomerID;references:ID" json:"addresses,omitempty"`
	Notes             []*CustomerNote    `gorm:"foreignKey:CustomerID;references:ID" json:"notes,omitempty"`
}

func (m *Customer) TableName() string {
	return CustomerTable
}

func (m *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(CustomerPrefix)
	}
	return
}

var CustomerSchema = struct {
	ID                schema.Field
	BusinessID        schema.Field
	Name              schema.Field
	CountryCode       schema.Field
	Gender            schema.Field
	Email             schema.Field
	PhoneNumber       schema.Field
	PhoneCode         schema.Field
	TikTokUsername    schema.Field
	InstagramUsername schema.Field
	FacebookUsername  schema.Field
	XUsername         schema.Field
	SnapchatUsername  schema.Field
	WhatsappNumber    schema.Field
	JoinedAt          schema.Field
	CreatedAt         schema.Field
	UpdatedAt         schema.Field
	DeletedAt         schema.Field
}{
	ID:                schema.NewField("id", "id"),
	BusinessID:        schema.NewField("business_id", "businessId"),
	Name:              schema.NewField("name", "name"),
	CountryCode:       schema.NewField("country_code", "countryCode"),
	Gender:            schema.NewField("gender", "gender"),
	Email:             schema.NewField("email", "email"),
	PhoneNumber:       schema.NewField("phone_number", "phoneNumber"),
	PhoneCode:         schema.NewField("phone_code", "phoneCode"),
	TikTokUsername:    schema.NewField("tiktok_username", "tiktokUsername"),
	InstagramUsername: schema.NewField("instagram_username", "instagramUsername"),
	FacebookUsername:  schema.NewField("facebook_username", "facebookUsername"),
	XUsername:         schema.NewField("x_username", "xUsername"),
	SnapchatUsername:  schema.NewField("snapchat_username", "snapchatUsername"),
	WhatsappNumber:    schema.NewField("whatsapp_number", "whatsappNumber"),
	JoinedAt:          schema.NewField("joined_at", "joinedAt"),
	CreatedAt:         schema.NewField("created_at", "createdAt"),
	UpdatedAt:         schema.NewField("updated_at", "updatedAt"),
	DeletedAt:         schema.NewField("deleted_at", "deletedAt"),
}

type CreateCustomerRequest struct {
	Name              string         `json:"name" binding:"required"`
	Gender            CustomerGender `json:"gender" binding:"omitempty,oneof=male female other"`
	CountryCode       string         `json:"countryCode" binding:"required,len=2"`
	Email             string         `json:"email" binding:"required,email"`
	PhoneNumber       string         `json:"phoneNumber" binding:"omitempty"`
	PhoneCode         string         `json:"phoneCode" binding:"omitempty"`
	TikTokUsername    string         `json:"tiktokUsername" binding:"omitempty"`
	InstagramUsername string         `json:"instagramUsername" binding:"omitempty"`
	FacebookUsername  string         `json:"facebookUsername" binding:"omitempty"`
	XUsername         string         `json:"xUsername" binding:"omitempty"`
	SnapchatUsername  string         `json:"snapchatUsername" binding:"omitempty"`
	JoinedAt          time.Time      `json:"joinedAt" binding:"omitempty"`
	WhatsappNumber    string         `json:"whatsappNumber" binding:"omitempty"`
}

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

/* Address Model */
//------------------*/

const (
	CustomerAddressTable  = "customer_addresses"
	CustomerAddressStruct = "CustomerAddress"
	CustomerAddressPrefix = "addr"
)

type CustomerAddress struct {
	gorm.Model
	ID          string         `gorm:"column:id;primaryKey;type:text" json:"id"`
	CustomerID  string         `gorm:"column:customer_id;type:text;not null;index" json:"customerId"`
	Customer    *Customer      `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	CountryCode string         `gorm:"column:country_code;type:text; not null" json:"countryCode"`
	State       string         `gorm:"column:state;type:text; not null" json:"state"`
	City        string         `gorm:"column:city;type:text; not null" json:"city"`
	Street      sql.NullString `gorm:"column:street;type:text" json:"street"`
	PhoneCode   string         `gorm:"column:phone_code;type:text; not null" json:"phoneCode"`
	PhoneNumber string         `gorm:"column:phone_number;type:text; not null" json:"phoneNumber"`
	ZipCode     sql.NullString `gorm:"column:zip_code;type:text" json:"zipCode"`
}

func (m *CustomerAddress) TableName() string {
	return CustomerAddressTable
}

func (m *CustomerAddress) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(CustomerAddressPrefix)
	}
	return
}

var CustomerAddressSchema = struct {
	ID          schema.Field
	CustomerID  schema.Field
	Street      schema.Field
	City        schema.Field
	State       schema.Field
	CountryCode schema.Field
	PhoneCode   schema.Field
	PhoneNumber schema.Field
	ZipCode     schema.Field
	CreatedAt   schema.Field
	UpdatedAt   schema.Field
	DeletedAt   schema.Field
}{
	ID:          schema.NewField("id", "id"),
	CustomerID:  schema.NewField("customer_id", "customerId"),
	Street:      schema.NewField("street", "street"),
	City:        schema.NewField("city", "city"),
	State:       schema.NewField("state", "state"),
	CountryCode: schema.NewField("country_code", "countryCode"),
	PhoneCode:   schema.NewField("phone_code", "phoneCode"),
	PhoneNumber: schema.NewField("phone_number", "phoneNumber"),
	ZipCode:     schema.NewField("zip_code", "zipCode"),
	CreatedAt:   schema.NewField("created_at", "createdAt"),
	UpdatedAt:   schema.NewField("updated_at", "updatedAt"),
	DeletedAt:   schema.NewField("deleted_at", "deletedAt"),
}

type CreateCustomerAddressRequest struct {
	CountryCode string `json:"countryCode" binding:"required,len=2"`
	State       string `json:"state" binding:"required"`
	City        string `json:"city" binding:"required"`
	PhoneCode   string `json:"phoneCode" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Street      string `json:"street" binding:"omitempty"`
	ZipCode     string `json:"zipCode" binding:"omitempty"`
}

type UpdateCustomerAddressRequest struct {
	Street      string `json:"street" binding:"omitempty"`
	City        string `json:"city" binding:"omitempty"`
	State       string `json:"state" binding:"omitempty"`
	CountryCode string `json:"countryCode" binding:"omitempty,len=2"`
	PhoneCode   string `json:"phoneCode" binding:"omitempty"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty"`
	ZipCode     string `json:"zipCode" binding:"omitempty"`
}

/* Customer Note Model */
//-----------------------*/

const (
	CustomerNoteTable  = "customer_notes"
	CustomerNoteStruct = "CustomerNote"
	CustomerNotePrefix = "cnote"
)

type CustomerNote struct {
	gorm.Model
	ID         string    `gorm:"column:id;primaryKey;type:text" json:"id"`
	CustomerID string    `gorm:"column:customer_id;type:text;not null;index" json:"customerId"`
	Customer   *Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	Content    string    `gorm:"column:content;type:text;not null" json:"content"`
}

func (m *CustomerNote) TableName() string {
	return CustomerNoteTable
}

func (m *CustomerNote) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = id.KsuidWithPrefix(CustomerNotePrefix)
	}
	return
}

var CustomerNoteSchema = struct {
	ID         schema.Field
	CustomerID schema.Field
	Content    schema.Field
	CreatedAt  schema.Field
	UpdatedAt  schema.Field
	DeletedAt  schema.Field
}{
	ID:         schema.NewField("id", "id"),
	CustomerID: schema.NewField("customer_id", "customerId"),
	Content:    schema.NewField("content", "content"),
	CreatedAt:  schema.NewField("created_at", "createdAt"),
	UpdatedAt:  schema.NewField("updated_at", "updatedAt"),
	DeletedAt:  schema.NewField("deleted_at", "deletedAt"),
}

type CreateCustomerNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateCustomerNoteRequest struct {
	Content string `json:"content" binding:"omitempty"`
}
