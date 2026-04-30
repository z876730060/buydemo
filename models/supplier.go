package models

import "time"

type Supplier struct {
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
