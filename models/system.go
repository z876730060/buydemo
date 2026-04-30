package models

import "time"

type SystemSetting struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"uniqueIndex;size:128;not null"`
	Value     string    `json:"value" gorm:"type:text"`
	Desc      string    `json:"desc" gorm:"size:256"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BackupRecord struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FileName   string    `json:"file_name" gorm:"size:256;not null"`
	FilePath   string    `json:"file_path" gorm:"size:512;not null"`
	FileSize   int64     `json:"file_size" gorm:"default:0"`
	Status     string    `json:"status" gorm:"size:32;default:success"` // success / failed
	Remark     string    `json:"remark" gorm:"type:text"`
	CreatedBy  uint      `json:"created_by" gorm:"not null"`
	Creator    User      `json:"creator" gorm:"foreignKey:CreatedBy"`
	CreatedAt  time.Time `json:"created_at"`
}
