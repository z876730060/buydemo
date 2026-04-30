package models

import "time"

type Contract struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	ContractNo    string     `json:"contract_no" gorm:"uniqueIndex;size:64;not null"`
	Title         string     `json:"title" gorm:"size:256;not null"`
	Type          string     `json:"type" gorm:"size:32;not null"` // purchase / sales
	PartyBType    string     `json:"party_b_type" gorm:"size:32;not null"` // supplier / customer
	PartyBID      uint       `json:"party_b_id" gorm:"not null"`
	PartyBName    string     `json:"party_b_name" gorm:"size:128"`
	PartyBContact string     `json:"party_b_contact" gorm:"size:64"`
	Content       string     `json:"content" gorm:"type:text"`
	TotalAmount   float64    `json:"total_amount" gorm:"type:decimal(12,2);default:0"`
	Status        string     `json:"status" gorm:"size:32;default:draft"` // draft / sent / signed_by_b / completed / cancelled
	SignToken     string     `json:"sign_token" gorm:"uniqueIndex;size:64"` // UUID for public signing link
	SignedByA     string     `json:"signed_by_a" gorm:"type:text"`        // 甲方签名 base64
	SignedByB     string     `json:"signed_by_b" gorm:"type:text"`        // 乙方签名 base64
	SignedAtA     *time.Time `json:"signed_at_a"`
	SignedAtB     *time.Time `json:"signed_at_b"`
	CreatedBy     uint       `json:"created_by" gorm:"not null"`
	Creator       User       `json:"creator" gorm:"foreignKey:CreatedBy"`
	Remark        string     `json:"remark" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
