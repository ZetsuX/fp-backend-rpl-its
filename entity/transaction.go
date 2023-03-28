package entity

import "fp-rpl/common"

type Transaction struct {
	common.Model
	Code 	   string   `json:"code" binding:"required"`
	TotalPrice float64  `json:"total_price" binding:"required"`
	Timestamp  string   `json:"timestamp" binding:"required"`
	Spots      []Spot   `json:"spot,omitempty" binding:"required"`
	UserID     uint64   `gorm:"foreignKey" json:"user_id" binding:"required"`
	User       *User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	SessionID  uint64   `gorm:"foreignKey" json:"session_id" binding:"required"`
	Session    *Session `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"session,omitempty"`
}
