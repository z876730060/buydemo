package models

import "time"

type Product struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Code      string    `json:"code" gorm:"uniqueIndex;size:64;not null"`
	Name      string    `json:"name" gorm:"size:128;not null"`
	Category  string    `json:"category" gorm:"size:64"`
	Unit      string    `json:"unit" gorm:"size:16"`
	Spec      string    `json:"spec" gorm:"size:128"`
	Price     float64   `json:"price" gorm:"type:decimal(12,2);default:0"`
	Status    int       `json:"status" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
