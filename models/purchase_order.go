package models

import "time"

type PurchaseOrder struct {
	ID           uint                 `json:"id" gorm:"primaryKey"`
	OrderNo      string               `json:"order_no" gorm:"uniqueIndex;size:64;not null"`
	SupplierID   uint                 `json:"supplier_id" gorm:"not null"`
	Supplier     Supplier             `json:"supplier" gorm:"foreignKey:SupplierID"`
	WarehouseID  uint                 `json:"warehouse_id" gorm:"not null;index"`
	Warehouse    Warehouse            `json:"warehouse" gorm:"foreignKey:WarehouseID"`
	TotalAmount  float64              `json:"total_amount" gorm:"type:decimal(12,2);default:0"`
	Status       string               `json:"status" gorm:"size:32;default:draft"`
	Remark       string               `json:"remark" gorm:"type:text"`
	CreatedBy    uint                 `json:"created_by" gorm:"not null"`
	Creator      User                 `json:"creator" gorm:"foreignKey:CreatedBy"`
	ApprovedAt   *time.Time           `json:"approved_at"`
	ReceivedAt   *time.Time           `json:"received_at"`
	Items        []PurchaseOrderItem  `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

type PurchaseOrderItem struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	OrderID      uint           `json:"order_id" gorm:"not null"`
	ProductID    uint           `json:"product_id" gorm:"not null"`
	Product      Product        `json:"product" gorm:"foreignKey:ProductID"`
	Quantity     float64        `json:"quantity" gorm:"type:decimal(12,2);not null"`
	UnitPrice    float64        `json:"unit_price" gorm:"type:decimal(12,2);not null"`
	Amount       float64        `json:"amount" gorm:"type:decimal(12,2);default:0"`
	ReceivedQty  float64        `json:"received_qty" gorm:"type:decimal(12,2);default:0"`
}
