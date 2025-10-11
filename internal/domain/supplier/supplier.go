package supplier

import (
	"github.com/abdelrahman146/kyora/internal/utils"
	"gorm.io/gorm"
)

const (
	SupplierTable  = "suppliers"
	SupplierAlias  = "sup"
	SupplierStruct = "Supplier"
)

type Supplier struct {
	gorm.Model
	ID      string `gorm:"column:id;primaryKey;type:text" json:"id"`
	StoreID string `gorm:"column:store_id;type:text;not null;index" json:"storeId"`
	Name    string `gorm:"column:name;type:text;not null" json:"name"`
	Contact string `gorm:"column:contact;type:text" json:"contact,omitempty"`
	Email   string `gorm:"column:email;type:text" json:"email,omitempty"`
	Phone   string `gorm:"column:phone;type:text" json:"phone,omitempty"`
	Website string `gorm:"column:website;type:text" json:"website,omitempty"`
}

func (m *Supplier) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(SupplierAlias)
	}
	return
}

type CreateSupplierRequest struct {
	Name    string `form:"name" json:"name" binding:"required"`
	Contact string `form:"contact" json:"contact" binding:"omitempty"`
	Email   string `form:"email" json:"email" binding:"omitempty,email"`
	Phone   string `form:"phone" json:"phone" binding:"omitempty"`
	Website string `form:"website" json:"website" binding:"omitempty,url"`
}

type UpdateSupplierRequest struct {
	Name    string `form:"name" json:"name" binding:"omitempty"`
	Contact string `form:"contact" json:"contact" binding:"omitempty"`
	Email   string `form:"email" json:"email" binding:"omitempty,email"`
	Phone   string `form:"phone" json:"phone" binding:"omitempty"`
	Website string `form:"website" json:"website" binding:"omitempty,url"`
}
