package models

import "time"

type OperationLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Username  string    `json:"username" gorm:"size:64"`
	Action    string    `json:"action" gorm:"size:64;not null"`   // create / update / delete / login / export / approve / receive / deliver / pay
	Target    string    `json:"target" gorm:"size:64;not null"`   // supplier / product / purchase_order / sales_order / user / expense / finance
	TargetID  uint      `json:"target_id"`                         // 操作对象ID
	Detail    string    `json:"detail" gorm:"type:text"`           // 操作详情/描述
	IP        string    `json:"ip" gorm:"size:64"`
	CreatedAt time.Time `json:"created_at"`
}
