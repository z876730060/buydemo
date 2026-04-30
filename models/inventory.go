package models

import "time"

type Inventory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id" gorm:"uniqueIndex;not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  float64   `json:"quantity" gorm:"type:decimal(12,2);default:0"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InventoryLog struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ProductID     uint      `json:"product_id" gorm:"not null"`
	Product       Product   `json:"product" gorm:"foreignKey:ProductID"`
	Type          string    `json:"type" gorm:"size:32;not null"` // in / out / adjust
	Quantity      float64   `json:"quantity" gorm:"type:decimal(12,2);not null"`
	Balance       float64   `json:"balance" gorm:"type:decimal(12,2);default:0"`
	ReferenceType string    `json:"reference_type" gorm:"size:32"` // purchase / adjust
	ReferenceID   uint      `json:"reference_id"`
	Remark        string    `json:"remark" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at"`
}
