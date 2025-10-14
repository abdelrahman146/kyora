package order

import (
	"github.com/abdelrahman146/kyora/internal/domain/inventory"
	"github.com/abdelrahman146/kyora/internal/utils"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	OrderItemTable  = "order_items"
	OrderItemAlias  = "oitm"
	OrderItemStruct = "OrderItem"
)

type OrderItem struct {
	gorm.Model
	ID        string             `gorm:"column:id;primaryKey;type:text" json:"id"`
	OrderID   string             `gorm:"column:order_id;type:text;not null;index" json:"orderId"`
	Order     *Order             `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
	ProductID string             `gorm:"column:product_id;type:text;not null;index" json:"productId"`
	Product   *inventory.Product `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	Quantity  int                `gorm:"column:quantity;type:int;not null;default:1" json:"quantity"`
	Currency  string             `gorm:"column:currency;type:text;not null" json:"currency"`
	UnitPrice decimal.Decimal    `gorm:"column:unit_price;type:numeric;not null;default:0" json:"unitPrice"`
	UnitCost  decimal.Decimal    `gorm:"column:unit_cost;type:numeric;not null;default:0" json:"unitCost"`
	TotalCost decimal.Decimal    `gorm:"column:total_cost;type:numeric;not null;default:0" json:"totalCost"`
	Total     decimal.Decimal    `gorm:"column:total;type:numeric;not null;default:0" json:"total"`
}

func (m *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = utils.ID.NewUlidWithPrefix(OrderItemAlias)
	}
	return
}

type CreateOrderItemRequest struct {
	ProductID string          `json:"productId" binding:"required"`
	Quantity  int             `json:"quantity" binding:"required,min=1"`
	UnitPrice decimal.Decimal `json:"unitPrice" binding:"required,gt=0"`
	UnitCost  decimal.Decimal `json:"unitCost" binding:"omitempty,gte=0"`
}

type UpdateOrderItemRequest struct {
	Quantity  int             `json:"quantity" binding:"omitempty,min=1"`
	UnitPrice decimal.Decimal `json:"unitPrice" binding:"omitempty,gt=0"`
	UnitCost  decimal.Decimal `json:"unitCost" binding:"omitempty,gte=0"`
}
