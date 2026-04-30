package models

import "time"

type Customer struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Code          string    `json:"code" gorm:"uniqueIndex;size:64;not null"`
	Name          string    `json:"name" gorm:"size:128;not null"`
	ContactPerson string    `json:"contact_person" gorm:"size:64"`
	Phone         string    `json:"phone" gorm:"size:32"`
	Address       string    `json:"address" gorm:"type:text"`
	Status        int       `json:"status" gorm:"default:1"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SalesOrder struct {
	ID          uint                `json:"id" gorm:"primaryKey"`
	OrderNo     string              `json:"order_no" gorm:"uniqueIndex;size:64;not null"`
	CustomerID  uint                `json:"customer_id" gorm:"not null"`
	Customer    Customer            `json:"customer" gorm:"foreignKey:CustomerID"`
	TotalAmount float64             `json:"total_amount" gorm:"type:decimal(12,2);default:0"`
	Status      string              `json:"status" gorm:"size:32;default:draft"`
	Remark      string              `json:"remark" gorm:"type:text"`
	CreatedBy   uint                `json:"created_by" gorm:"not null"`
	Creator     User                `json:"creator" gorm:"foreignKey:CreatedBy"`
	ApprovedAt  *time.Time          `json:"approved_at"`
	DeliveredAt *time.Time          `json:"delivered_at"`
	Items       []SalesOrderItem    `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type SalesOrderItem struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	OrderID       uint      `json:"order_id" gorm:"not null"`
	ProductID     uint      `json:"product_id" gorm:"not null"`
	Product       Product   `json:"product" gorm:"foreignKey:ProductID"`
	Quantity      float64   `json:"quantity" gorm:"type:decimal(12,2);not null"`
	UnitPrice     float64   `json:"unit_price" gorm:"type:decimal(12,2);not null"`
	Amount        float64   `json:"amount" gorm:"type:decimal(12,2);default:0"`
	DeliveredQty  float64   `json:"delivered_qty" gorm:"type:decimal(12,2);default:0"`
}
