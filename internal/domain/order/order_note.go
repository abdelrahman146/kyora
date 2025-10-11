package order

import (
	"github.com/abdelrahman146/kyora/internal/utils"
	"gorm.io/gorm"
)

const (
	OrderNoteTable  = "order_notes"
	OrderNoteAlias  = "onot"
	OrderNoteStruct = "OrderNote"
)

type OrderNote struct {
	gorm.Model
	ID      string `gorm:"column:id;primaryKey;type:text" json:"id"`
	OrderID string `gorm:"column:order_id;type:text;not null;index" json:"orderId"`
	Order   *Order `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
	Note    string `gorm:"column:note;type:text;not null" json:"note"`
}

type CreateOrderNoteRequest struct {
	OrderID string `json:"orderId" binding:"required"`
	Note    string `json:"note" binding:"required"`
}

type UpdateOrderNoteRequest struct {
	Note string `json:"note" binding:"omitempty"`
}

func (m *OrderNote) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(OrderNoteAlias)
	}
	return
}
