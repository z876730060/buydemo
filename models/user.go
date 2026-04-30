package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:64;not null"`
	Password  string    `json:"-" gorm:"size:256;not null"`
	RealName  string    `json:"real_name" gorm:"size:64"`
	Role      string    `json:"role" gorm:"size:32;default:operator"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
