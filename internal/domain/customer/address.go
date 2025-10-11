package customer

import (
	"github.com/abdelrahman146/kyora/internal/utils"
	"gorm.io/gorm"
)

const (
	AddressTable  = "customer_addresses"
	AddressAlias  = "addr"
	AddressStruct = "Address"
)

type Address struct {
	gorm.Model
	ID          string    `gorm:"column:id;primaryKey;type:text" json:"id"`
	CustomerID  string    `gorm:"column:customer_id;type:text;not null;index" json:"customerId"`
	Customer    *Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	Street      string    `gorm:"column:street;type:text" json:"street"`
	City        string    `gorm:"column:city;type:text" json:"city"`
	State       string    `gorm:"column:state;type:text" json:"state"`
	CountryCode string    `gorm:"column:country_code;type:text" json:"countryCode"`
	Phone       string    `gorm:"column:phone;type:text" json:"phone"`
	ZipCode     string    `gorm:"column:zip_code;type:text" json:"zipCode"`
}

func (m *Address) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(AddressAlias)
	}
	return
}

type CreateAddressRequest struct {
	Street      string `json:"street" binding:"required"`
	City        string `json:"city" binding:"required"`
	State       string `json:"state" binding:"omitempty"`
	CountryCode string `json:"countryCode" binding:"required,len=2"`
	Phone       string `json:"phone" binding:"omitempty"`
	ZipCode     string `json:"zipCode" binding:"omitempty"`
}

type UpdateAddressRequest struct {
	Street      string `json:"street" binding:"omitempty"`
	City        string `json:"city" binding:"omitempty"`
	State       string `json:"state" binding:"omitempty"`
	CountryCode string `json:"countryCode" binding:"omitempty,len=2"`
	Phone       string `json:"phone" binding:"omitempty"`
	ZipCode     string `json:"zipCode" binding:"omitempty"`
}
