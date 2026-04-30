package models

import "time"

// 应付账款（采购入库自动生成）
type AccountPayable struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	OrderNo       string     `json:"order_no" gorm:"size:64;not null"`
	PurchaseOrderID uint     `json:"purchase_order_id" gorm:"not null"`
	SupplierID    uint       `json:"supplier_id" gorm:"not null"`
	Supplier      Supplier   `json:"supplier" gorm:"foreignKey:SupplierID"`
	TotalAmount   float64    `json:"total_amount" gorm:"type:decimal(12,2);default:0"`
	PaidAmount    float64    `json:"paid_amount" gorm:"type:decimal(12,2);default:0"`
	DueAmount     float64    `json:"due_amount" gorm:"type:decimal(12,2);default:0"`
	Status        string     `json:"status" gorm:"size:32;default:pending"` // pending / partial / paid
	DueDate       *time.Time `json:"due_date"`
	PaidAt        *time.Time `json:"paid_at"`
	Remark        string     `json:"remark" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// 应收账款（销售出库自动生成）
type AccountReceivable struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	OrderNo       string     `json:"order_no" gorm:"size:64;not null"`
	SalesOrderID  uint       `json:"sales_order_id" gorm:"not null"`
	CustomerID    uint       `json:"customer_id" gorm:"not null"`
	Customer      Customer   `json:"customer" gorm:"foreignKey:CustomerID"`
	TotalAmount   float64    `json:"total_amount" gorm:"type:decimal(12,2);default:0"`
	ReceivedAmount float64   `json:"received_amount" gorm:"type:decimal(12,2);default:0"`
	DueAmount     float64    `json:"due_amount" gorm:"type:decimal(12,2);default:0"`
	Status        string     `json:"status" gorm:"size:32;default:pending"` // pending / partial / received
	DueDate       *time.Time `json:"due_date"`
	ReceivedAt    *time.Time `json:"received_at"`
	Remark        string     `json:"remark" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// 费用记录
type Expense struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Category    string    `json:"category" gorm:"size:64;not null"` // 办公费/运输费/税费/其他
	Amount      float64   `json:"amount" gorm:"type:decimal(12,2);not null"`
	Description string    `json:"description" gorm:"type:text"`
	OccurredAt  time.Time `json:"occurred_at"`
	CreatedBy   uint      `json:"created_by" gorm:"not null"`
	Creator     User      `json:"creator" gorm:"foreignKey:CreatedBy"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 付款记录
type PaymentRecord struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	PayableID       *uint     `json:"payable_id"`                          // 关联应付账款（可选）
	ReceivableID    *uint     `json:"receivable_id"`                       // 关联应收账款（可选）
	Type            string    `json:"type" gorm:"size:32;not null"`        // pay (付款) / receive (收款)
	Amount          float64   `json:"amount" gorm:"type:decimal(12,2);not null"`
	PaymentMethod   string    `json:"payment_method" gorm:"size:32"`       // cash / bank / wechat / alipay
	ReferenceNo     string    `json:"reference_no" gorm:"size:128"`        // 凭证号/银行流水号
	Remark          string    `json:"remark" gorm:"type:text"`
	OperatedBy      uint      `json:"operated_by" gorm:"not null"`
	Operator        User      `json:"operator" gorm:"foreignKey:OperatedBy"`
	OperatedAt      time.Time `json:"operated_at"`
	CreatedAt       time.Time `json:"created_at"`
}
