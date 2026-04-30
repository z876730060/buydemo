package models

import "time"

type Company struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Code        string    `json:"code" gorm:"uniqueIndex;size:64;not null"`
	Name        string    `json:"name" gorm:"size:128;not null"`
	Contact     string    `json:"contact" gorm:"size:64"`
	Phone       string    `json:"phone" gorm:"size:32"`
	Email       string    `json:"email" gorm:"size:128"`
	Address     string    `json:"address" gorm:"type:text"`
	TaxNo       string    `json:"tax_no" gorm:"size:64"`
	Status      int       `json:"status" gorm:"default:1"` // 1: 启用, 0: 禁用
	Remark      string    `json:"remark" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Warehouse struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Code        string    `json:"code" gorm:"uniqueIndex;size:64;not null"`
	Name        string    `json:"name" gorm:"size:128;not null"`
	CompanyID   uint      `json:"company_id" gorm:"not null;index"`
	Company     Company   `json:"company" gorm:"foreignKey:CompanyID"`
	Address     string    `json:"address" gorm:"type:text"`
	Contact     string    `json:"contact" gorm:"size:64"`
	Phone       string    `json:"phone" gorm:"size:32"`
	IsDefault   bool      `json:"is_default" gorm:"default:false"` // 默认仓库
	Status      int       `json:"status" gorm:"default:1"`         // 1: 启用, 0: 禁用
	Remark      string    `json:"remark" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
