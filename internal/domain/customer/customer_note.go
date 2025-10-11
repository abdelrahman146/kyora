package customer

import (
	"github.com/abdelrahman146/kyora/internal/utils"
	"gorm.io/gorm"
)

const (
	CustomerNoteTable  = "customer_notes"
	CustomerNoteAlias  = "cnot"
	CustomerNoteStruct = "CustomerNote"
)

type CustomerNote struct {
	ID         string    `gorm:"column:id;primaryKey;type:text" json:"id"`
	CustomerID string    `gorm:"column:customer_id;type:text;not null;index" json:"customerId"`
	Customer   *Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	Note       string    `gorm:"column:note;type:text;not null" json:"note"`
	CreatedAt  int64     `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  int64     `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (m *CustomerNote) BeforeCreate(tx gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(CustomerNoteAlias)
	}
	return
}

type CreateCustomerNoteRequest struct {
	CustomerID string `json:"customerId" binding:"required"`
	Note       string `json:"note" binding:"required"`
}

type UpdateCustomerNoteRequest struct {
	Note string `json:"note" binding:"omitempty"`
}
