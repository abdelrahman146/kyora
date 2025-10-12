package owner

import (
	"github.com/abdelrahman146/kyora/internal/domain/store"
	"github.com/abdelrahman146/kyora/internal/utils"
	"gorm.io/gorm"
)

const (
	OwnerTable  = "owners"
	OwnerStruct = "Owner"
	OwnerAlias  = "own"
)

type Owner struct {
	gorm.Model
	ID        string       `json:"id" gorm:"column:id;primaryKey;type:text"`
	StoreID   string       `json:"storeId" gorm:"column:store_id;type:text;not null;index"`
	Store     *store.Store `json:"store,omitempty" gorm:"foreignKey:StoreID;references:ID"`
	FirstName string       `json:"firstName" gorm:"column:first_name;type:text;not null"`
	LastName  string       `json:"lastName" gorm:"column:last_name;type:text;not null"`
}

func (m *Owner) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(OwnerAlias)
	}
	return
}

type CreateOwnerRequest struct {
	FirstName string `form:"firstName" json:"firstName" binding:"required"`
	LastName  string `form:"lastName" json:"lastName" binding:"required"`
}

type UpdateOwnerRequest struct {
	FirstName string `form:"firstName" json:"firstName" binding:"omitempty"`
	LastName  string `form:"lastName" json:"lastName" binding:"omitempty"`
}
