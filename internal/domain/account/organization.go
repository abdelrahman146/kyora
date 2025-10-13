package account

import (
	"github.com/abdelrahman146/kyora/internal/utils"

	"gorm.io/gorm"
)

const (
	OrganizationTable  = "organizations"
	OrganizationStruct = "Organization"
	OrganizationAlias  = "org"
)

type Organization struct {
	gorm.Model
	ID   string `gorm:"column:id;primaryKey;type:text" json:"id"`
	Name string `gorm:"column:name;type:text;not null" json:"name"`
}

func (m *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(OrganizationAlias)
	}
	return
}

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}
